package delivery

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/delivery"
	"telegramshop_backend/pkg/deliverykind"
)

type Service interface {
	GetAll(ctx context.Context) ([]models.DeliveryMethod, error)
	GetByID(ctx context.Context, id string) (models.DeliveryMethod, error)
	Create(ctx context.Context, m models.DeliveryMethod) (models.DeliveryMethod, error)
	Update(ctx context.Context, id string, m models.DeliveryMethod) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo delivery.Repository
}

func NewService(repo delivery.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]models.DeliveryMethod, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (models.DeliveryMethod, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, m models.DeliveryMethod) (models.DeliveryMethod, error) {
	m.Kind = deliverykind.Normalize(m.Kind)
	return s.repo.Create(ctx, m)
}

func (s *service) Update(ctx context.Context, id string, m models.DeliveryMethod) error {
	m.Kind = deliverykind.Normalize(m.Kind)
	return s.repo.Update(ctx, id, m)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
