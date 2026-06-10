package basket

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/basket"
	"telegramshop_backend/pkg/logger"
)

type Service interface {
	GetUserBasket(ctx context.Context, telegramID int64) ([]models.BasketItem, error)
	AddToBasket(ctx context.Context, input models.CreateBasketItem) error
	UpdateBasketItem(ctx context.Context, input models.CreateBasketItem) error
	RemoveFromBasket(ctx context.Context, telegramID int64, productID string) error
}

type service struct {
	repo basket.Repository
}

func NewService(repo basket.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserBasket(ctx context.Context, telegramID int64) ([]models.BasketItem, error) {
	return s.repo.GetUserBasket(ctx, telegramID)
}

func (s *service) AddToBasket(ctx context.Context, input models.CreateBasketItem) error {
	logger.Infof("[AddToBasket] tg=%d product=%s qty=%d", input.UserID, input.ProductID, input.Quantity)
	return s.repo.CreateOrUpdateBasketItem(ctx, input.UserID, input.ProductID, input.Quantity)
}

func (s *service) UpdateBasketItem(ctx context.Context, input models.CreateBasketItem) error {
	return s.repo.CreateOrUpdateBasketItem(ctx, input.UserID, input.ProductID, input.Quantity)
}

func (s *service) RemoveFromBasket(ctx context.Context, telegramID int64, productID string) error {
	return s.repo.DeleteBasketItem(ctx, telegramID, productID)
}
