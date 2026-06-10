package payments

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAll(ctx context.Context) ([]models.PaymentMethod, error)
	GetByID(ctx context.Context, id string) (models.PaymentMethod, error)
	Create(ctx context.Context, m models.PaymentMethod) (models.PaymentMethod, error)
	Update(ctx context.Context, id string, m models.PaymentMethod) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	err := r.db.SelectContext(ctx, &methods, `
		SELECT id, name, description, is_active, created_at
		FROM payment_methods
		ORDER BY name`)
	return methods, err
}

func (r *repository) GetByID(ctx context.Context, id string) (models.PaymentMethod, error) {
	var m models.PaymentMethod
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, is_active, created_at
		FROM payment_methods WHERE id = $1`, id).Scan(
		&m.ID, &m.Name, &m.Description, &m.IsActive, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return models.PaymentMethod{}, sql.ErrNoRows
	}
	return m, err
}

func (r *repository) Create(ctx context.Context, m models.PaymentMethod) (models.PaymentMethod, error) {
	var out models.PaymentMethod
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO payment_methods (name, description, created_at)
		VALUES ($1, $2, now())
		RETURNING id, name, description, is_active, created_at`,
		m.Name, m.Description).Scan(
		&out.ID, &out.Name, &out.Description, &out.IsActive, &out.CreatedAt)
	return out, err
}

func (r *repository) Update(ctx context.Context, id string, m models.PaymentMethod) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE payment_methods SET name = $1, description = $2, is_active = $3
		WHERE id = $4`, m.Name, m.Description, m.IsActive, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM payment_methods WHERE id = $1`, id)
	return err
}
