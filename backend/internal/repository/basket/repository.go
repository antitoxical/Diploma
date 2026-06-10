package basket

import (
	"context"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserBasket(ctx context.Context, telegramID int64) ([]models.BasketItem, error)
	ClearUserBasket(ctx context.Context, telegramID int64) error
	CreateOrUpdateBasketItem(ctx context.Context, telegramID int64, productID string, quantity int) error
	DeleteBasketItem(ctx context.Context, telegramID int64, productID string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserBasket(ctx context.Context, telegramID int64) ([]models.BasketItem, error) {
	query := `
		SELECT bi.id, bi.user_id, bi.product_id, bi.quantity, bi.created_at
		FROM basket_items bi
		INNER JOIN users u ON u.id = bi.user_id
		WHERE u.telegram_id = $1`

	var items []models.BasketItem
	err := r.db.SelectContext(ctx, &items, query, telegramID)
	return items, err
}

func (r *repository) ClearUserBasket(ctx context.Context, telegramID int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM basket_items bi USING users u
		WHERE bi.user_id = u.id AND u.telegram_id = $1`, telegramID)
	return err
}

func (r *repository) CreateOrUpdateBasketItem(ctx context.Context, telegramID int64, productID string, quantity int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO basket_items (user_id, product_id, quantity, created_at)
		SELECT u.id, $2::uuid, $3, now()
		FROM users u WHERE u.telegram_id = $1
		ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = EXCLUDED.quantity`,
		telegramID, productID, quantity)
	return err
}

func (r *repository) DeleteBasketItem(ctx context.Context, telegramID int64, productID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM basket_items bi USING users u
		WHERE bi.user_id = u.id AND u.telegram_id = $1 AND bi.product_id = $2::uuid`,
		telegramID, productID)
	return err
}
