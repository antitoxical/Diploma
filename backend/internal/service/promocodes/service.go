package promocodes

import (
	"context"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/promocodes"
	"telegramshop_backend/internal/repository/users"
	"telegramshop_backend/pkg/logger"
	"telegramshop_backend/pkg/telegram"
)

type Service interface {
	GetAll(ctx context.Context) ([]models.Promocode, error)
	GetByID(ctx context.Context, id string) (models.Promocode, error)
	ValidateCode(ctx context.Context, code string) (models.Promocode, error)
	Create(ctx context.Context, input models.Promocode) (models.Promocode, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo     promocodes.Repository
	userRepo users.Repository
	bot      *telegram.BotClient
}

func NewService(repo promocodes.Repository, userRepo users.Repository, bot *telegram.BotClient) Service {
	return &service{repo: repo, userRepo: userRepo, bot: bot}
}

func (s *service) GetAll(ctx context.Context) ([]models.Promocode, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (models.Promocode, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ValidateCode(ctx context.Context, code string) (models.Promocode, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *service) Create(ctx context.Context, input models.Promocode) (models.Promocode, error) {
	promo, err := s.repo.Create(ctx, input)
	if err != nil {
		return models.Promocode{}, err
	}
	if s.bot != nil && s.userRepo != nil {
		go s.broadcastPromocode(promo)
	}
	return promo, nil
}

func (s *service) broadcastPromocode(promo models.Promocode) {
	ctx := context.Background()
	ids, err := s.userRepo.ListTelegramIDs(ctx)
	if err != nil {
		logger.Errorf("[promocodes] list telegram ids: %v", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	msg := telegram.FormatPromocodeMessage(promo.Code, promo.DiscountPercent)
	s.bot.BroadcastMessage(ids, msg, true)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
