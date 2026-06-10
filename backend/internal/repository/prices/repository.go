package prices

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateTier(ctx context.Context, tier models.ProductPriceTier) (models.ProductPriceTier, error)
	GetTierByID(ctx context.Context, id string) (models.ProductPriceTier, error)
	GetTiersByProductID(ctx context.Context, productID string) ([]models.ProductPriceTier, error)
	UpdateTier(ctx context.Context, id string, input models.UpdatePriceTierInput) error
	DeleteTier(ctx context.Context, id string) error
	DeleteTiersByProductID(ctx context.Context, productID string) error
	UpdateTierMinQuantity(ctx context.Context, id string, newMin int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTier(ctx context.Context, tier models.ProductPriceTier) (models.ProductPriceTier, error) {
	query := `
		INSERT INTO product_price_tiers (product_id, min_quantity, price_kopeks, created_at)
		VALUES ($1, $2, $3, now())
		RETURNING id, product_id, min_quantity, price_kopeks, created_at`

	var out models.ProductPriceTier
	err := r.db.QueryRowContext(ctx, query, tier.ProductID, tier.MinQuantity, tier.PriceKopeks).Scan(
		&out.ID, &out.ProductID, &out.MinQuantity, &out.PriceKopeks, &out.CreatedAt,
	)
	return out, err
}

func (r *repository) GetTierByID(ctx context.Context, id string) (models.ProductPriceTier, error) {
	query := `SELECT id, product_id, min_quantity, price_kopeks, created_at FROM product_price_tiers WHERE id = $1`
	var t models.ProductPriceTier
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.ProductID, &t.MinQuantity, &t.PriceKopeks, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return models.ProductPriceTier{}, sql.ErrNoRows
	}
	return t, err
}

func (r *repository) GetTiersByProductID(ctx context.Context, productID string) ([]models.ProductPriceTier, error) {
	query := `
		SELECT id, product_id, min_quantity, price_kopeks, created_at
		FROM product_price_tiers
		WHERE product_id = $1
		ORDER BY min_quantity`

	var tiers []models.ProductPriceTier
	err := r.db.SelectContext(ctx, &tiers, query, productID)
	return tiers, err
}

func (r *repository) UpdateTier(ctx context.Context, id string, input models.UpdatePriceTierInput) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE product_price_tiers SET min_quantity = $1, price_kopeks = $2 WHERE id = $3`,
		input.MinQuantity, input.PriceKopeks, id)
	return err
}

func (r *repository) DeleteTier(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_price_tiers WHERE id = $1`, id)
	return err
}

func (r *repository) DeleteTiersByProductID(ctx context.Context, productID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_price_tiers WHERE product_id = $1`, productID)
	return err
}

func (r *repository) UpdateTierMinQuantity(ctx context.Context, id string, newMin int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE product_price_tiers SET min_quantity = $1 WHERE id = $2`, newMin, id)
	return err
}
