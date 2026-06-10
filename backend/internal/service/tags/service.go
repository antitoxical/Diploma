package tags

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/tags"
)

type Service interface {
	Create(ctx context.Context, input models.CreateTagInput) (models.Tag, error)
	GetByID(ctx context.Context, id string) (models.Tag, error)
	GetAll(ctx context.Context) ([]models.Tag, error)
	Update(ctx context.Context, id string, input models.CreateTagInput) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo tags.Repository
}

func NewService(repo tags.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input models.CreateTagInput) (models.Tag, error) {
	return s.repo.Create(ctx, input)
}

func (s *service) GetByID(ctx context.Context, id string) (models.Tag, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]models.Tag, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) Update(ctx context.Context, id string, input models.CreateTagInput) error {
	return s.repo.Update(ctx, id, input)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
