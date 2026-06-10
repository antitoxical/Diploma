package favorites

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/favorites"
	"telegramshop_backend/pkg/logger"
)

type Service interface {
	GetUserFavorites(ctx context.Context, telegramID int64) ([]models.Favorite, error)
	AddToFavorites(ctx context.Context, input models.CreateFavorite) error
	RemoveFromFavorites(ctx context.Context, telegramID int64, productID string) error
}

type service struct {
	repo favorites.Repository
}

func NewService(repo favorites.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserFavorites(ctx context.Context, telegramID int64) ([]models.Favorite, error) {
	return s.repo.GetUserFavorites(ctx, telegramID)
}

func (s *service) AddToFavorites(ctx context.Context, input models.CreateFavorite) error {
	logger.Infof("[AddToFavorites] tg=%d product=%s", input.UserID, input.ProductID)
	return s.repo.CreateFavorite(ctx, input.UserID, input.ProductID)
}

func (s *service) RemoveFromFavorites(ctx context.Context, telegramID int64, productID string) error {
	return s.repo.DeleteFavorite(ctx, telegramID, productID)
}
