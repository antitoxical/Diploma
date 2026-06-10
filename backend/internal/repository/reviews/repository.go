package reviews

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	UpsertReview(ctx context.Context, userInternalID string, productID string, rating int, comment *string) (models.Review, error)
	DeleteReview(ctx context.Context, telegramID int64, productID string) error
	GetAll(ctx context.Context) ([]models.Review, error)
	GetByProduct(ctx context.Context, productID string) ([]models.Review, error)
	GetUserReview(ctx context.Context, telegramID int64, productID string) (models.Review, error)
	GetRatingStat(ctx context.Context, productID string) (models.ProductRatingStat, error)
	ModerateReview(ctx context.Context, reviewID string, approved bool) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) UpsertReview(ctx context.Context, userInternalID string, productID string, rating int, comment *string) (models.Review, error) {
	query := `
		INSERT INTO reviews (user_id, product_id, rating, comment, is_moderated, created_at)
		VALUES ($1, $2, $3, $4, false, now())
		ON CONFLICT (user_id, product_id) DO UPDATE SET
			rating = EXCLUDED.rating,
			comment = EXCLUDED.comment,
			is_moderated = false
		RETURNING id, user_id, product_id, rating, comment, is_moderated, created_at`

	var rev models.Review
	err := r.db.QueryRowContext(ctx, query, userInternalID, productID, rating, comment).Scan(
		&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating, &rev.Comment, &rev.IsModerated, &rev.CreatedAt,
	)
	return rev, err
}

func (r *repository) DeleteReview(ctx context.Context, telegramID int64, productID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM reviews r USING users u
		WHERE r.user_id = u.id AND u.telegram_id = $1 AND r.product_id = $2::uuid`,
		telegramID, productID)
	return err
}

func (r *repository) GetAll(ctx context.Context) ([]models.Review, error) {
	query := `
		SELECT id, user_id, product_id, rating, comment, is_moderated, created_at
		FROM reviews ORDER BY created_at DESC`
	var list []models.Review
	err := r.db.SelectContext(ctx, &list, query)
	return list, err
}

func (r *repository) GetByProduct(ctx context.Context, productID string) ([]models.Review, error) {
	query := `
		SELECT id, user_id, product_id, rating, comment, is_moderated, created_at
		FROM reviews WHERE product_id = $1 ORDER BY created_at DESC`

	var list []models.Review
	err := r.db.SelectContext(ctx, &list, query, productID)
	return list, err
}

func (r *repository) GetUserReview(ctx context.Context, telegramID int64, productID string) (models.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.product_id, r.rating, r.comment, r.is_moderated, r.created_at
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		WHERE u.telegram_id = $1 AND r.product_id = $2::uuid`

	var rev models.Review
	err := r.db.QueryRowContext(ctx, query, telegramID, productID).Scan(
		&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating, &rev.Comment, &rev.IsModerated, &rev.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Review{}, err
	}
	return rev, err
}

func (r *repository) GetRatingStat(ctx context.Context, productID string) (models.ProductRatingStat, error) {
	query := `
		SELECT product_id, avg_rating, review_count
		FROM product_rating_stats
		WHERE product_id = $1::uuid`

	var s models.ProductRatingStat
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&s.ProductID, &s.AvgRating, &s.ReviewCount)
	if err == sql.ErrNoRows {
		return models.ProductRatingStat{ProductID: productID, AvgRating: 0, ReviewCount: 0}, nil
	}
	return s, err
}

func (r *repository) ModerateReview(ctx context.Context, reviewID string, approved bool) error {
	if !approved {
		_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, reviewID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `UPDATE reviews SET is_moderated = true WHERE id = $1`, reviewID)
	return err
}
