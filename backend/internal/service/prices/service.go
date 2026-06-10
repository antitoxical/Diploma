package prices

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/prices"
	"telegramshop_backend/pkg/logger"
)

type Service interface {
	CreateTier(ctx context.Context, input models.ProductPriceTier) (models.ProductPriceTier, error)
	GetTierByID(ctx context.Context, id string) (models.ProductPriceTier, error)
	GetTiersByProductID(ctx context.Context, productID string) ([]models.ProductPriceTier, error)
	UpdateTier(ctx context.Context, id string, input models.UpdatePriceTierInput) error
	DeleteTier(ctx context.Context, id string) error
	DeleteTiersByProductID(ctx context.Context, productID string) error
	UpdateTierMinQuantity(ctx context.Context, id string, newMin int) error
}

type service struct {
	repo prices.Repository
}

func NewService(repo prices.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTier(ctx context.Context, input models.ProductPriceTier) (models.ProductPriceTier, error) {
	return s.repo.CreateTier(ctx, input)
}

func (s *service) GetTierByID(ctx context.Context, id string) (models.ProductPriceTier, error) {
	return s.repo.GetTierByID(ctx, id)
}

func (s *service) GetTiersByProductID(ctx context.Context, productID string) ([]models.ProductPriceTier, error) {
	return s.repo.GetTiersByProductID(ctx, productID)
}

func (s *service) UpdateTier(ctx context.Context, id string, input models.UpdatePriceTierInput) error {
	return s.repo.UpdateTier(ctx, id, input)
}

func (s *service) DeleteTier(ctx context.Context, id string) error {
	return s.repo.DeleteTier(ctx, id)
}

func (s *service) DeleteTiersByProductID(ctx context.Context, productID string) error {
	return s.repo.DeleteTiersByProductID(ctx, productID)
}

func (s *service) UpdateTierMinQuantity(ctx context.Context, id string, newMin int) error {
	if err := s.repo.UpdateTierMinQuantity(ctx, id, newMin); err != nil {
		logger.Errorf("[UpdateTierMinQuantity] %v", err)
		return err
	}
	return nil
}
