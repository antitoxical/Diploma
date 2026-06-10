package firms

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateFirm(ctx context.Context, firm models.Firm) (models.Firm, error)
	GetFirmByID(ctx context.Context, id string) (models.Firm, error)
	GetAllFirms(ctx context.Context) ([]models.Firm, error)
	UpdateFirm(ctx context.Context, id string, input models.UpdateFirmInput) error
	DeleteFirm(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateFirm(ctx context.Context, firm models.Firm) (models.Firm, error) {
	query := `INSERT INTO firms (name, created_at) VALUES ($1, now()) RETURNING id, name, created_at`
	var out models.Firm
	err := r.db.QueryRowContext(ctx, query, firm.Name).Scan(&out.ID, &out.Name, &out.CreatedAt)
	return out, err
}

func (r *repository) GetFirmByID(ctx context.Context, id string) (models.Firm, error) {
	query := `SELECT id, name, created_at FROM firms WHERE id = $1`
	var firm models.Firm
	err := r.db.QueryRowContext(ctx, query, id).Scan(&firm.ID, &firm.Name, &firm.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Firm{}, sql.ErrNoRows
	}
	return firm, err
}

func (r *repository) GetAllFirms(ctx context.Context) ([]models.Firm, error) {
	var firms []models.Firm
	err := r.db.SelectContext(ctx, &firms, `SELECT id, name, created_at FROM firms ORDER BY name`)
	return firms, err
}

func (r *repository) UpdateFirm(ctx context.Context, id string, input models.UpdateFirmInput) error {
	_, err := r.db.ExecContext(ctx, `UPDATE firms SET name = $1 WHERE id = $2`, input.Name, id)
	return err
}

func (r *repository) DeleteFirm(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM firms WHERE id = $1`, id)
	return err
}
