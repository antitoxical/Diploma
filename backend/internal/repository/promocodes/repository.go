package promocodes

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(ctx context.Context) ([]models.Promocode, error)
	GetByID(ctx context.Context, id string) (models.Promocode, error)
	GetByCode(ctx context.Context, code string) (models.Promocode, error)
	Create(ctx context.Context, input models.Promocode) (models.Promocode, error)
	Delete(ctx context.Context, id string) error
	IncrementUses(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]models.Promocode, error) {
	var promos []models.Promocode
	err := r.db.SelectContext(ctx, &promos, `
		SELECT id, code, discount_percent, max_uses, current_uses, is_active, expires_at, created_at
		FROM promocodes
		ORDER BY created_at DESC`)
	return promos, err
}

func (r *repository) GetByID(ctx context.Context, id string) (models.Promocode, error) {
	var p models.Promocode
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, discount_percent, max_uses, current_uses, is_active, expires_at, created_at
		FROM promocodes WHERE id = $1`, id).Scan(
		&p.ID, &p.Code, &p.DiscountPercent, &p.MaxUses, &p.CurrentUses, &p.IsActive, &p.ExpiresAt, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Promocode{}, sql.ErrNoRows
	}
	return p, err
}

func (r *repository) GetByCode(ctx context.Context, code string) (models.Promocode, error) {
	var p models.Promocode
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, discount_percent, max_uses, current_uses, is_active, expires_at, created_at
		FROM promocodes WHERE code = $1`, code).Scan(
		&p.ID, &p.Code, &p.DiscountPercent, &p.MaxUses, &p.CurrentUses, &p.IsActive, &p.ExpiresAt, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Promocode{}, sql.ErrNoRows
	}
	return p, err
}

func (r *repository) Create(ctx context.Context, input models.Promocode) (models.Promocode, error) {
	var p models.Promocode
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO promocodes (code, discount_percent, max_uses, expires_at, created_at)
		VALUES ($1, $2, $3, $4, now())
		RETURNING id, code, discount_percent, max_uses, current_uses, is_active, expires_at, created_at`,
		input.Code, input.DiscountPercent, input.MaxUses, input.ExpiresAt).Scan(
		&p.ID, &p.Code, &p.DiscountPercent, &p.MaxUses, &p.CurrentUses, &p.IsActive, &p.ExpiresAt, &p.CreatedAt)
	return p, err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM promocodes WHERE id = $1`, id)
	return err
}

func (r *repository) IncrementUses(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE promocodes SET current_uses = current_uses + 1 WHERE id = $1`, id)
	return err
}
