package notifications

import (
	"context"
	"fmt"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/notifications"
	"telegramshop_backend/internal/repository/users"
)

type Service interface {
	Create(ctx context.Context, input models.CreateNotification) (models.Notification, error)
	GetByUser(ctx context.Context, telegramID int64) ([]models.Notification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, telegramID int64) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo     notifications.Repository
	userRepo users.Repository
}

func NewService(repo notifications.Repository, userRepo users.Repository) Service {
	return &service{repo: repo, userRepo: userRepo}
}

func (s *service) Create(ctx context.Context, input models.CreateNotification) (models.Notification, error) {
	userUUID, err := s.userRepo.GetInternalIDByTelegramID(ctx, input.UserID)
	if err != nil {
		return models.Notification{}, fmt.Errorf("user lookup: %w", err)
	}
	return s.repo.Create(ctx, userUUID, input.Title, input.Message)
}

func (s *service) GetByUser(ctx context.Context, telegramID int64) ([]models.Notification, error) {
	return s.repo.GetByUser(ctx, telegramID)
}

func (s *service) MarkRead(ctx context.Context, id string) error {
	return s.repo.MarkRead(ctx, id)
}

func (s *service) MarkAllRead(ctx context.Context, telegramID int64) error {
	return s.repo.MarkAllRead(ctx, telegramID)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
