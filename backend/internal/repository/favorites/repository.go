package favorites

import (
	"context"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserFavorites(ctx context.Context, telegramID int64) ([]models.Favorite, error)
	CreateFavorite(ctx context.Context, telegramID int64, productID string) error
	DeleteFavorite(ctx context.Context, telegramID int64, productID string) error
	ListTelegramIDsByProductID(ctx context.Context, productID string) ([]int64, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserFavorites(ctx context.Context, telegramID int64) ([]models.Favorite, error) {
	query := `
		SELECT f.id, f.user_id, f.product_id, f.created_at
		FROM favorites f
		INNER JOIN users u ON u.id = f.user_id
		WHERE u.telegram_id = $1`

	var favorites []models.Favorite
	err := r.db.SelectContext(ctx, &favorites, query, telegramID)
	return favorites, err
}

func (r *repository) CreateFavorite(ctx context.Context, telegramID int64, productID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO favorites (user_id, product_id, created_at)
		SELECT u.id, $2::uuid, now()
		FROM users u WHERE u.telegram_id = $1
		ON CONFLICT (user_id, product_id) DO NOTHING`,
		telegramID, productID)
	return err
}

func (r *repository) DeleteFavorite(ctx context.Context, telegramID int64, productID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM favorites f USING users u
		WHERE f.user_id = u.id AND u.telegram_id = $1 AND f.product_id = $2::uuid`,
		telegramID, productID)
	return err
}

func (r *repository) ListTelegramIDsByProductID(ctx context.Context, productID string) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids, `
		SELECT DISTINCT u.telegram_id
		FROM favorites f
		INNER JOIN users u ON u.id = f.user_id
		WHERE f.product_id = $1::uuid
		  AND u.is_active = true
		  AND u.telegram_id IS NOT NULL
		  AND u.telegram_id != 0`, productID)
	return ids, err
}
