package reviews

import (
	"context"
	"database/sql"
	"errors"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/reviews"
	"telegramshop_backend/internal/repository/users"
)

type Service interface {
	UpsertReview(ctx context.Context, input models.UpsertReviewInput) (models.Review, error)
	DeleteReview(ctx context.Context, telegramID int64, productID string) error
	GetAll(ctx context.Context) ([]models.Review, error)
	GetByProduct(ctx context.Context, productID string) ([]models.Review, error)
	GetUserReview(ctx context.Context, telegramID int64, productID string) (models.Review, error)
	GetRatingStat(ctx context.Context, productID string) (models.ProductRatingStat, error)
	ModerateReview(ctx context.Context, reviewID string, approved bool) error
}

type service struct {
	repo     reviews.Repository
	userRepo users.Repository
}

func NewService(repo reviews.Repository, userRepo users.Repository) Service {
	return &service{repo: repo, userRepo: userRepo}
}

func (s *service) UpsertReview(ctx context.Context, input models.UpsertReviewInput) (models.Review, error) {
	if input.Rating < 1 || input.Rating > 5 {
		return models.Review{}, errors.New("rating must be 1..5")
	}
	uid, err := s.userRepo.GetInternalIDByTelegramID(ctx, input.TelegramID)
	if err != nil {
		return models.Review{}, err
	}
	var c *string
	if input.Comment != "" {
		c = &input.Comment
	}
	return s.repo.UpsertReview(ctx, uid, input.ProductID, input.Rating, c)
}

func (s *service) DeleteReview(ctx context.Context, telegramID int64, productID string) error {
	return s.repo.DeleteReview(ctx, telegramID, productID)
}

func (s *service) GetAll(ctx context.Context) ([]models.Review, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByProduct(ctx context.Context, productID string) ([]models.Review, error) {
	return s.repo.GetByProduct(ctx, productID)
}

func (s *service) GetUserReview(ctx context.Context, telegramID int64, productID string) (models.Review, error) {
	r, err := s.repo.GetUserReview(ctx, telegramID, productID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Review{}, err
	}
	return r, err
}

func (s *service) GetRatingStat(ctx context.Context, productID string) (models.ProductRatingStat, error) {
	return s.repo.GetRatingStat(ctx, productID)
}

func (s *service) ModerateReview(ctx context.Context, reviewID string, approved bool) error {
	return s.repo.ModerateReview(ctx, reviewID, approved)
}
