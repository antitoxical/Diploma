package delivery

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(ctx context.Context) ([]models.DeliveryMethod, error)
	GetByID(ctx context.Context, id string) (models.DeliveryMethod, error)
	Create(ctx context.Context, m models.DeliveryMethod) (models.DeliveryMethod, error)
	Update(ctx context.Context, id string, m models.DeliveryMethod) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]models.DeliveryMethod, error) {
	var methods []models.DeliveryMethod
	err := r.db.SelectContext(ctx, &methods, `
		SELECT id, name, kind, price_kopeks, description, is_active, created_at
		FROM delivery_methods
		ORDER BY name`)
	return methods, err
}

func (r *repository) GetByID(ctx context.Context, id string) (models.DeliveryMethod, error) {
	var m models.DeliveryMethod
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, kind, price_kopeks, description, is_active, created_at
		FROM delivery_methods WHERE id = $1`, id).Scan(
		&m.ID, &m.Name, &m.Kind, &m.PriceKopeks, &m.Description, &m.IsActive, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return models.DeliveryMethod{}, sql.ErrNoRows
	}
	return m, err
}

func (r *repository) Create(ctx context.Context, m models.DeliveryMethod) (models.DeliveryMethod, error) {
	var out models.DeliveryMethod
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO delivery_methods (name, kind, price_kopeks, description, created_at)
		VALUES ($1, $2, $3, $4, now())
		RETURNING id, name, kind, price_kopeks, description, is_active, created_at`,
		m.Name, m.Kind, m.PriceKopeks, m.Description).Scan(
		&out.ID, &out.Name, &out.Kind, &out.PriceKopeks, &out.Description, &out.IsActive, &out.CreatedAt)
	return out, err
}

func (r *repository) Update(ctx context.Context, id string, m models.DeliveryMethod) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE delivery_methods SET name = $1, kind = $2, price_kopeks = $3, description = $4, is_active = $5
		WHERE id = $6`, m.Name, m.Kind, m.PriceKopeks, m.Description, m.IsActive, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM delivery_methods WHERE id = $1`, id)
	return err
}
