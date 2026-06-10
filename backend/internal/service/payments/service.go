package payments

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/payments"
)

type Service interface {
	GetAll(ctx context.Context) ([]models.PaymentMethod, error)
	GetByID(ctx context.Context, id string) (models.PaymentMethod, error)
	Create(ctx context.Context, m models.PaymentMethod) (models.PaymentMethod, error)
	Update(ctx context.Context, id string, m models.PaymentMethod) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo payments.Repository
}

func NewService(repo payments.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]models.PaymentMethod, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (models.PaymentMethod, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, m models.PaymentMethod) (models.PaymentMethod, error) {
	return s.repo.Create(ctx, m)
}

func (s *service) Update(ctx context.Context, id string, m models.PaymentMethod) error {
	return s.repo.Update(ctx, id, m)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
