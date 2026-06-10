package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/orders"
	"telegramshop_backend/internal/repository/products"
	"telegramshop_backend/internal/repository/promocodes"
	"telegramshop_backend/internal/repository/users"
	notificationsService "telegramshop_backend/internal/service/notifications"
	"telegramshop_backend/pkg/logger"
	"telegramshop_backend/pkg/telegram"
	"telegramshop_backend/pkg/validate"

	"github.com/google/uuid"
)

var validStatuses = map[string]bool{
	"created":    true,
	"paid":       true,
	"delivering": true,
	"completed":  true,
	"cancelled":  true,
}

type Service interface {
	GetAll(ctx context.Context) ([]models.OrderWithProducts, error)
	CreateOrder(ctx context.Context, input models.CreateOrder) (models.OrderWithProducts, error)
	GetOrderByID(ctx context.Context, id string) (models.OrderWithProducts, error)
	GetUserOrders(ctx context.Context, telegramID int64) ([]models.OrderWithProducts, error)
	UpdateStatus(ctx context.Context, id string, input models.UpdateOrderStatus) error
	GetStatusHistory(ctx context.Context, orderID string) ([]models.OrderStatusHistory, error)
}

type service struct {
	repo          orders.Repository
	userRepo      users.Repository
	productRepo   products.Repository
	promocodeRepo promocodes.Repository
	notifications notificationsService.Service
	bot           *telegram.BotClient
}

func NewService(
	repo orders.Repository,
	userRepo users.Repository,
	productRepo products.Repository,
	promocodeRepo promocodes.Repository,
	notifications notificationsService.Service,
	bot *telegram.BotClient,
) Service {
	return &service{
		repo:          repo,
		userRepo:      userRepo,
		productRepo:   productRepo,
		promocodeRepo: promocodeRepo,
		notifications: notifications,
		bot:           bot,
	}
}

func (s *service) pushInAppNotification(ctx context.Context, telegramID int64, title, message string) {
	if s.notifications == nil || telegramID == 0 {
		return
	}
	_, err := s.notifications.Create(ctx, models.CreateNotification{
		UserID:  telegramID,
		Title:   title,
		Message: message,
	})
	if err != nil {
		logger.Errorf("[orders] in-app notification: %v", err)
	}
}

func (s *service) GetAll(ctx context.Context) ([]models.OrderWithProducts, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) CreateOrder(ctx context.Context, input models.CreateOrder) (models.OrderWithProducts, error) {
	if len(input.Items) == 0 {
		return models.OrderWithProducts{}, errors.New("empty order items")
	}
	if input.DeliveryMethodID == "" || input.PaymentMethodID == "" {
		return models.OrderWithProducts{}, errors.New("delivery and payment methods are required")
	}
	if _, err := uuid.Parse(input.DeliveryMethodID); err != nil {
		return models.OrderWithProducts{}, errors.New("invalid delivery method id")
	}
	if _, err := uuid.Parse(input.PaymentMethodID); err != nil {
		return models.OrderWithProducts{}, errors.New("invalid payment method id")
	}
	if err := validate.Address(input.DeliveryAddress, true); err != nil {
		return models.OrderWithProducts{}, err
	}

	userUUID, err := s.userRepo.GetInternalIDByTelegramID(ctx, input.UserID)
	if err != nil {
		return models.OrderWithProducts{}, fmt.Errorf("user: %w", err)
	}

	var lines []orders.OrderLineInput
	var totalKopeks int64
	var notifyItems []telegram.OrderLineNotify

	for _, it := range input.Items {
		if _, err := uuid.Parse(it.ProductID); err != nil {
			return models.OrderWithProducts{}, fmt.Errorf("invalid product id: %w", err)
		}
		if it.Quantity <= 0 {
			return models.OrderWithProducts{}, errors.New("invalid quantity")
		}

		p, err := s.productRepo.GetProductByID(ctx, it.ProductID)
		if err != nil {
			return models.OrderWithProducts{}, err
		}
		if !p.IsActive {
			return models.OrderWithProducts{}, fmt.Errorf("product %s inactive", it.ProductID)
		}
		if p.StockQuantity < it.Quantity {
			return models.OrderWithProducts{}, fmt.Errorf("insufficient stock for product %s", it.ProductID)
		}

		unitKopeks, err := s.productRepo.GetProductPriceForQuantity(ctx, it.ProductID, it.Quantity)
		if err != nil {
			return models.OrderWithProducts{}, err
		}

		lineTotal := unitKopeks * int64(it.Quantity)
		totalKopeks += lineTotal

		lines = append(lines, orders.OrderLineInput{
			ProductID:             it.ProductID,
			Quantity:              it.Quantity,
			PriceAtPurchaseKopeks: unitKopeks,
		})
		notifyItems = append(notifyItems, telegram.OrderLineNotify{
			Name:       p.Name,
			Quantity:   it.Quantity,
			UnitKopeks: unitKopeks,
		})
	}

	var dmID, pmID, pcID *string

	if input.DeliveryMethodID != "" {
		dmID = &input.DeliveryMethodID
	}
	if input.PaymentMethodID != "" {
		pmID = &input.PaymentMethodID
	}

	if input.PromoCode != "" {
		promo, err := s.promocodeRepo.GetByCode(ctx, input.PromoCode)
		if err == nil && promo.IsActive {
			validExpiry := promo.ExpiresAt == nil || promo.ExpiresAt.After(time.Now())
			validUses := promo.MaxUses == nil || promo.CurrentUses < *promo.MaxUses
			if validExpiry && validUses {
				pcID = &promo.ID
				discount := totalKopeks * int64(promo.DiscountPercent) / 100
				totalKopeks -= discount
				if totalKopeks < 0 {
					totalKopeks = 0
				}
				_ = s.promocodeRepo.IncrementUses(ctx, promo.ID)
			}
		}
	}

	order, err := s.repo.CreateOrder(ctx, input, userUUID, totalKopeks, lines, dmID, pmID, pcID)
	if err != nil {
		return models.OrderWithProducts{}, err
	}

	deliveryName, paymentName := s.repo.LookupMethodNames(ctx, dmID, pmID)
	notify := telegram.OrderCreatedNotify{
		OrderID:      order.ID,
		TotalKopeks:  totalKopeks,
		Items:        notifyItems,
		DeliveryName: deliveryName,
		PaymentName:  paymentName,
		Address:      input.DeliveryAddress,
	}

	title, body := telegram.FormatOrderCreatedInAppDetailed(notify)
	s.pushInAppNotification(ctx, input.UserID, title, body)

	if s.bot != nil && input.UserID != 0 {
		msg := telegram.FormatOrderCreatedMessageDetailed(notify)
		if sendErr := s.bot.SendMessageWithOrderButton(input.UserID, msg); sendErr != nil {
			logger.Errorf("[CreateOrder] failed to send telegram notification: %v", sendErr)
		}
	}

	if s.bot != nil && s.bot.HasAdmin() {
		adminMsg := telegram.FormatAdminNewOrderMessage(order.ID, totalKopeks)
		if sendErr := s.bot.SendAdminNewOrderNotification(order.ID, adminMsg); sendErr != nil {
			logger.Errorf("[CreateOrder] admin notification: %v", sendErr)
		}
	}

	return order, nil
}

func (s *service) GetOrderByID(ctx context.Context, id string) (models.OrderWithProducts, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		logger.Errorf("[GetOrderByID] %v", err)
		return models.OrderWithProducts{}, err
	}
	return order, nil
}

func (s *service) GetUserOrders(ctx context.Context, telegramID int64) ([]models.OrderWithProducts, error) {
	return s.repo.GetUserOrders(ctx, telegramID)
}

func (s *service) UpdateStatus(ctx context.Context, id string, input models.UpdateOrderStatus) error {
	if !validStatuses[input.Status] {
		return fmt.Errorf("invalid order status: %s", input.Status)
	}

	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get order for notification: %w", err)
	}

	var comment *string
	if input.Comment != "" {
		comment = &input.Comment
	}

	if err := s.repo.UpdateStatus(ctx, id, input.Status, comment); err != nil {
		return err
	}

	s.notifyOrderStatus(ctx, order, input.Status, input.Comment)

	return nil
}

func (s *service) notifyOrderStatus(ctx context.Context, order models.OrderWithProducts, status, comment string) {
	title, body := telegram.FormatOrderStatusInApp(order.ID, status, comment)
	s.pushInAppNotification(ctx, order.TelegramID, title, body)

	if s.bot == nil || order.TelegramID == 0 {
		return
	}

	switch status {
	case "cancelled":
		msg := telegram.FormatOrderCancelledMessage(order.ID, comment)
		if sendErr := s.bot.SendMessageWithOrderButton(order.TelegramID, msg); sendErr != nil {
			logger.Errorf("[UpdateStatus] cancelled notification: %v", sendErr)
		}
	case "completed":
		productID := firstOrderProductID(order.Products)
		msg := telegram.FormatOrderCompletedReviewMessage(order.ID)
		if sendErr := s.bot.SendMessageWithReviewButton(order.TelegramID, order.ID, productID, msg); sendErr != nil {
			logger.Errorf("[UpdateStatus] completed notification: %v", sendErr)
		}
	default:
		msg := telegram.FormatOrderStatusMessage(order.ID, status)
		if sendErr := s.bot.SendMessageWithOrderButton(order.TelegramID, msg); sendErr != nil {
			logger.Errorf("[UpdateStatus] status notification: %v", sendErr)
		}
	}
}

func firstOrderProductID(items []models.OrderItemRow) string {
	if len(items) == 0 {
		return ""
	}
	return items[0].ProductID
}

func (s *service) GetStatusHistory(ctx context.Context, orderID string) ([]models.OrderStatusHistory, error) {
	_, err := s.repo.GetOrderByID(ctx, orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, err
	}
	return s.repo.GetStatusHistory(ctx, orderID)
}
