package firms

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/firms"
	"telegramshop_backend/pkg/logger"
)

type Service interface {
	CreateFirm(ctx context.Context, input models.Firm) (models.Firm, error)
	GetFirmByID(ctx context.Context, id string) (models.Firm, error)
	GetAllFirms(ctx context.Context) ([]models.Firm, error)
	UpdateFirm(ctx context.Context, id string, input models.UpdateFirmInput) error
	DeleteFirm(ctx context.Context, id string) error
}

type service struct {
	repo firms.Repository
}

func NewService(repo firms.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateFirm(ctx context.Context, input models.Firm) (models.Firm, error) {
	logger.Infof("[CreateFirm] Creating firm with name=%s", input.Name)
	return s.repo.CreateFirm(ctx, input)
}

func (s *service) GetFirmByID(ctx context.Context, id string) (models.Firm, error) {
	firm, err := s.repo.GetFirmByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("[GetFirmByID] Firm not found: id=%s", id)
		}
		return models.Firm{}, err
	}
	return firm, nil
}

func (s *service) GetAllFirms(ctx context.Context) ([]models.Firm, error) {
	return s.repo.GetAllFirms(ctx)
}

func (s *service) UpdateFirm(ctx context.Context, id string, input models.UpdateFirmInput) error {
	return s.repo.UpdateFirm(ctx, id, input)
}

func (s *service) DeleteFirm(ctx context.Context, id string) error {
	return s.repo.DeleteFirm(ctx, id)
}
