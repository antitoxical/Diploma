package products

import (
	"context"

	"telegramshop_backend/internal/models"
	favoritesRepo "telegramshop_backend/internal/repository/favorites"
	"telegramshop_backend/internal/repository/products"
	"telegramshop_backend/pkg/logger"
	"telegramshop_backend/pkg/telegram"
)

type Service interface {
	CreateProduct(ctx context.Context, input models.Product) (models.Product, error)
	GetProductByID(ctx context.Context, id string) (models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	SearchProducts(ctx context.Context, filter models.ProductFilter) (models.PaginatedProducts, error)
	UpdateProduct(ctx context.Context, id string, input models.UpdateProductInput) error
	DeleteProduct(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, productID string, stock int) error

	AddProductImage(ctx context.Context, productID string, imageURL string, isMain bool) (models.ProductImage, error)
	RemoveProductImage(ctx context.Context, imageID string) error
	SetProductImages(ctx context.Context, productID string, images []string) error

	AddProductTag(ctx context.Context, productID string, tagID string) error
	RemoveProductTag(ctx context.Context, productID string, tagID string) error
}

type service struct {
	repo    products.Repository
	favRepo favoritesRepo.Repository
	bot     *telegram.BotClient
}

func NewService(repo products.Repository, favRepo favoritesRepo.Repository, bot *telegram.BotClient) Service {
	return &service{repo: repo, favRepo: favRepo, bot: bot}
}

func (s *service) CreateProduct(ctx context.Context, input models.Product) (models.Product, error) {
	urls := input.ImageURLs
	input.ImageURLs = nil

	p, err := s.repo.CreateProduct(ctx, input)
	if err != nil {
		return models.Product{}, err
	}
	if len(urls) > 0 {
		if err := s.repo.SetProductImages(ctx, p.ID, urls); err != nil {
			return models.Product{}, err
		}
	}
	return s.repo.GetProductByID(ctx, p.ID)
}

func (s *service) GetProductByID(ctx context.Context, id string) (models.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *service) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAllProducts(ctx)
}

func (s *service) SearchProducts(ctx context.Context, filter models.ProductFilter) (models.PaginatedProducts, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	items, total, err := s.repo.SearchProducts(ctx, filter)
	if err != nil {
		return models.PaginatedProducts{}, err
	}
	if items == nil {
		items = []models.Product{}
	}

	totalPages := total / filter.Limit
	if total%filter.Limit != 0 {
		totalPages++
	}

	return models.PaginatedProducts{
		Products:   items,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *service) UpdateProduct(ctx context.Context, id string, input models.UpdateProductInput) error {
	imageURLs := input.ImageURLs

	if err := s.repo.UpdateProduct(ctx, id, input); err != nil {
		return err
	}
	if imageURLs != nil {
		return s.repo.SetProductImages(ctx, id, *imageURLs)
	}
	return nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *service) AddProductImage(ctx context.Context, productID string, imageURL string, isMain bool) (models.ProductImage, error) {
	return s.repo.AddProductImage(ctx, productID, imageURL, isMain)
}

func (s *service) RemoveProductImage(ctx context.Context, imageID string) error {
	return s.repo.RemoveProductImage(ctx, imageID)
}

func (s *service) SetProductImages(ctx context.Context, productID string, images []string) error {
	return s.repo.SetProductImages(ctx, productID, images)
}

func (s *service) UpdateStock(ctx context.Context, productID string, stock int) error {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	oldStock := product.StockQuantity

	logger.Infof("[UpdateStock] product=%s stock=%d", productID, stock)
	if err := s.repo.UpdateStock(ctx, productID, stock); err != nil {
		return err
	}

	if oldStock <= 0 && stock > 0 && s.bot != nil && s.favRepo != nil {
		go s.notifyBackInStock(productID, product.Name)
	}
	return nil
}

func (s *service) notifyBackInStock(productID, productName string) {
	ctx := context.Background()
	ids, err := s.favRepo.ListTelegramIDsByProductID(ctx, productID)
	if err != nil {
		logger.Errorf("[UpdateStock] favorites telegram ids: %v", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	msg := telegram.FormatBackInStockMessage(productName)
	for _, chatID := range ids {
		if err := s.bot.SendMessageWithProductButton(chatID, productID, msg); err != nil {
			logger.Errorf("[UpdateStock] back-in-stock notify %d: %v", chatID, err)
		}
	}
}

func (s *service) AddProductTag(ctx context.Context, productID string, tagID string) error {
	return s.repo.AddProductTag(ctx, productID, tagID)
}

func (s *service) RemoveProductTag(ctx context.Context, productID string, tagID string) error {
	return s.repo.RemoveProductTag(ctx, productID, tagID)
}
