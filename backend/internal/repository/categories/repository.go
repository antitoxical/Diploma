package categories

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
	GetCategoryByID(ctx context.Context, id string) (models.Category, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetRootCategories(ctx context.Context) ([]models.Category, error)
	GetSubcategories(ctx context.Context, parentID string) ([]models.Category, error)
	UpdateCategory(ctx context.Context, id string, category models.UpdateCategoryInput) error
	DeleteCategory(ctx context.Context, id string) error
	SetImage(ctx context.Context, id string, imageURL string) error
	RemoveImage(ctx context.Context, id string) error
}

const categoryCols = `id, name, parent_id, image_url, is_active, created_at`

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	query := `
		INSERT INTO categories (name, parent_id, image_url, created_at)
		VALUES ($1, $2, $3, now())
		RETURNING ` + categoryCols

	var out models.Category
	err := r.db.QueryRowContext(ctx, query, category.Name, category.ParentID, category.ImageURL).Scan(
		&out.ID, &out.Name, &out.ParentID, &out.ImageURL, &out.IsActive, &out.CreatedAt,
	)
	return out, err
}

func (r *repository) GetCategoryByID(ctx context.Context, id string) (models.Category, error) {
	query := `SELECT ` + categoryCols + ` FROM categories WHERE id = $1`

	var c models.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.ParentID, &c.ImageURL, &c.IsActive, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Category{}, err
	}
	return c, err
}

func (r *repository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	query := `SELECT ` + categoryCols + ` FROM categories WHERE is_active = true ORDER BY name`

	var categories []models.Category
	err := r.db.SelectContext(ctx, &categories, query)
	return categories, err
}

func (r *repository) GetRootCategories(ctx context.Context) ([]models.Category, error) {
	query := `SELECT ` + categoryCols + ` FROM categories WHERE parent_id IS NULL AND is_active = true ORDER BY name`

	var categories []models.Category
	err := r.db.SelectContext(ctx, &categories, query)
	return categories, err
}

func (r *repository) GetSubcategories(ctx context.Context, parentID string) ([]models.Category, error) {
	query := `SELECT ` + categoryCols + ` FROM categories WHERE parent_id = $1 AND is_active = true ORDER BY name`

	var categories []models.Category
	err := r.db.SelectContext(ctx, &categories, query, parentID)
	return categories, err
}

func (r *repository) UpdateCategory(ctx context.Context, id string, category models.UpdateCategoryInput) error {
	query := `UPDATE categories SET name = $1, parent_id = $2, image_url = $3, is_active = COALESCE($4, is_active) WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.ParentID, category.ImageURL, category.IsActive, id)
	return err
}

func (r *repository) DeleteCategory(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}

func (r *repository) SetImage(ctx context.Context, id string, imageURL string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE categories SET image_url = $1 WHERE id = $2`, imageURL, id)
	return err
}

func (r *repository) RemoveImage(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE categories SET image_url = NULL WHERE id = $1`, id)
	return err
}
