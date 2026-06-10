package tags

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, input models.CreateTagInput) (models.Tag, error)
	GetByID(ctx context.Context, id string) (models.Tag, error)
	GetAll(ctx context.Context) ([]models.Tag, error)
	Update(ctx context.Context, id string, input models.CreateTagInput) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input models.CreateTagInput) (models.Tag, error) {
	var tag models.Tag
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO tags (name, color, created_at)
		VALUES ($1, $2, now())
		RETURNING id, name, color, created_at`,
		input.Name, input.Color).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
	return tag, err
}

func (r *repository) GetByID(ctx context.Context, id string) (models.Tag, error) {
	var tag models.Tag
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, color, created_at FROM tags WHERE id = $1`, id).Scan(
		&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Tag{}, sql.ErrNoRows
	}
	return tag, err
}

func (r *repository) GetAll(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.SelectContext(ctx, &tags, `SELECT id, name, color, created_at FROM tags ORDER BY name`)
	return tags, err
}

func (r *repository) Update(ctx context.Context, id string, input models.CreateTagInput) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tags SET name = $1, color = $2 WHERE id = $3`,
		input.Name, input.Color, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id = $1`, id)
	return err
}
