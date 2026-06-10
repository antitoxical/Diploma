package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProductByID(ctx context.Context, id string) (models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	SearchProducts(ctx context.Context, filter models.ProductFilter) ([]models.Product, int, error)
	UpdateProduct(ctx context.Context, id string, product models.UpdateProductInput) error
	DeleteProduct(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, productID string, stock int) error
	GetProductPriceForQuantity(ctx context.Context, productID string, quantity int) (kopeks int64, err error)

	AddProductImage(ctx context.Context, productID string, imageURL string, isMain bool) (models.ProductImage, error)
	RemoveProductImage(ctx context.Context, imageID string) error
	SetProductImages(ctx context.Context, productID string, images []string) error
	GetProductImages(ctx context.Context, productID string) ([]models.ProductImage, error)

	GetProductTags(ctx context.Context, productID string) ([]models.Tag, error)
	AddProductTag(ctx context.Context, productID string, tagID string) error
	RemoveProductTag(ctx context.Context, productID string, tagID string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

const productCols = `id, category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active, created_at`

func scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}, p *models.Product) error {
	var firmID, desc sql.NullString
	err := scanner.Scan(
		&p.ID,
		&p.CategoryID,
		&firmID,
		&p.Name,
		&desc,
		&p.PriceKopeks,
		&p.StockQuantity,
		&p.IsActive,
		&p.CreatedAt,
	)
	if err != nil {
		return err
	}
	if firmID.Valid {
		s := firmID.String
		p.FirmID = &s
	}
	if desc.Valid {
		s := desc.String
		p.Description = &s
	}
	return nil
}

func (r *repository) loadImages(ctx context.Context, p *models.Product) error {
	imgs, err := r.GetProductImages(ctx, p.ID)
	if err != nil {
		return err
	}
	p.Images = imgs
	return nil
}

func (r *repository) loadTags(ctx context.Context, p *models.Product) error {
	tags, err := r.GetProductTags(ctx, p.ID)
	if err != nil {
		return err
	}
	p.Tags = tags
	return nil
}

func (r *repository) enrichProduct(ctx context.Context, p *models.Product) error {
	if err := r.loadImages(ctx, p); err != nil {
		return err
	}
	return r.loadTags(ctx, p)
}

func (r *repository) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	query := `
		INSERT INTO products (category_id, firm_id, name, description, price_kopeks, stock_quantity, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		RETURNING ` + productCols

	var out models.Product
	row := r.db.QueryRowxContext(ctx, query,
		product.CategoryID,
		product.FirmID,
		product.Name,
		product.Description,
		product.PriceKopeks,
		product.StockQuantity,
		product.IsActive,
	)
	if err := scanProduct(row, &out); err != nil {
		return models.Product{}, err
	}
	return out, nil
}

func (r *repository) GetProductByID(ctx context.Context, id string) (models.Product, error) {
	query := `SELECT ` + productCols + ` FROM products WHERE id = $1`

	row := r.db.QueryRowxContext(ctx, query, id)
	var p models.Product
	err := scanProduct(row, &p)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Product{}, err
	}
	if err != nil {
		return models.Product{}, err
	}
	if err := r.enrichProduct(ctx, &p); err != nil {
		return models.Product{}, err
	}
	return p, nil
}

func (r *repository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT ` + productCols + ` FROM products ORDER BY created_at DESC`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProduct(rows, &p); err != nil {
			return nil, err
		}
		if err := r.enrichProduct(ctx, &p); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *repository) SearchProducts(ctx context.Context, filter models.ProductFilter) ([]models.Product, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, "p.is_active = true")

	if filter.CategoryID != nil && *filter.CategoryID != "" {
		if filter.IncludeSubcategories {
			conditions = append(conditions, fmt.Sprintf(`p.category_id IN (
				WITH RECURSIVE cat_tree AS (
					SELECT id FROM categories WHERE id = $%d::uuid
					UNION ALL
					SELECT c.id FROM categories c INNER JOIN cat_tree ct ON c.parent_id = ct.id
				) SELECT id FROM cat_tree
			)`, argIdx))
		} else {
			conditions = append(conditions, fmt.Sprintf("p.category_id = $%d::uuid", argIdx))
		}
		args = append(args, *filter.CategoryID)
		argIdx++
	}

	if filter.FirmID != nil && *filter.FirmID != "" {
		conditions = append(conditions, fmt.Sprintf("p.firm_id = $%d::uuid", argIdx))
		args = append(args, *filter.FirmID)
		argIdx++
	}

	if filter.TagID != nil && *filter.TagID != "" {
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM product_tags pt WHERE pt.product_id = p.id AND pt.tag_id = $%d::uuid)", argIdx))
		args = append(args, *filter.TagID)
		argIdx++
	}

	if len(filter.TagIDs) > 0 {
		placeholders := make([]string, len(filter.TagIDs))
		for i, tid := range filter.TagIDs {
			placeholders[i] = fmt.Sprintf("$%d::uuid", argIdx)
			args = append(args, tid)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf(
			`(SELECT COUNT(DISTINCT pt.tag_id) FROM product_tags pt WHERE pt.product_id = p.id AND pt.tag_id IN (%s)) = %d`,
			strings.Join(placeholders, ","), len(filter.TagIDs)))
	}

	if filter.Search != nil && *filter.Search != "" {
		pattern := "%" + strings.ToLower(*filter.Search) + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(LOWER(p.name) LIKE $%d OR LOWER(COALESCE(p.description,'')) LIKE $%d)", argIdx, argIdx))
		args = append(args, pattern)
		argIdx++
	}

	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price_kopeks >= $%d", argIdx))
		args = append(args, *filter.MinPrice)
		argIdx++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price_kopeks <= $%d", argIdx))
		args = append(args, *filter.MaxPrice)
		argIdx++
	}

	if filter.InStock != nil && *filter.InStock {
		conditions = append(conditions, "p.stock_quantity > 0")
	}

	if filter.NewDays != nil && *filter.NewDays > 0 {
		conditions = append(conditions, fmt.Sprintf("p.created_at >= NOW() - INTERVAL '%d days'", *filter.NewDays))
	}

	if filter.HasReviews != nil && *filter.HasReviews {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM reviews rv WHERE rv.product_id = p.id AND rv.is_moderated = true)")
	}

	if filter.MinRating != nil && *filter.MinRating > 0 {
		conditions = append(conditions, fmt.Sprintf(
			`(SELECT COALESCE(AVG(rv.rating), 0) FROM reviews rv WHERE rv.product_id = p.id AND rv.is_moderated = true) >= $%d`, argIdx))
		args = append(args, *filter.MinRating)
		argIdx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	countQuery := "SELECT COUNT(*) FROM products p " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := resolveOrderBy(filter.SortBy)

	offset := (filter.Page - 1) * filter.Limit
	dataQuery := fmt.Sprintf(`
		SELECT %s FROM products p %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`, "p."+strings.ReplaceAll(productCols, ", ", ", p."), where, orderBy, argIdx, argIdx+1)
	args = append(args, filter.Limit, offset)

	rows, err := r.db.QueryxContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProduct(rows, &p); err != nil {
			return nil, 0, err
		}
		if err := r.enrichProduct(ctx, &p); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}
	return list, total, rows.Err()
}

func resolveOrderBy(sortBy string) string {
	switch sortBy {
	case "price_asc":
		return "p.price_kopeks ASC"
	case "price_desc":
		return "p.price_kopeks DESC"
	case "newest":
		return "p.created_at DESC"
	case "rating":
		return "(SELECT COALESCE(AVG(rv.rating), 0) FROM reviews rv WHERE rv.product_id = p.id AND rv.is_moderated = true) DESC"
	case "popular":
		return "(SELECT COUNT(*) FROM order_items oi WHERE oi.product_id = p.id) DESC"
	default:
		return "p.created_at DESC"
	}
}

func (r *repository) UpdateProduct(ctx context.Context, id string, product models.UpdateProductInput) error {
	query := `
		UPDATE products
		SET name = $1,
			category_id = $2,
			firm_id = $3,
			description = $4,
			price_kopeks = $5,
			stock_quantity = $6,
			is_active = COALESCE($7, is_active)
		WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.CategoryID,
		product.FirmID,
		product.Description,
		product.PriceKopeks,
		product.StockQuantity,
		product.IsActive,
		id,
	)
	return err
}

func (r *repository) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *repository) UpdateStock(ctx context.Context, productID string, stock int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE products SET stock_quantity = $1 WHERE id = $2`, stock, productID)
	return err
}

func (r *repository) GetProductPriceForQuantity(ctx context.Context, productID string, quantity int) (int64, error) {
	var tierK int64
	err := r.db.QueryRowContext(ctx, `
		SELECT price_kopeks FROM product_price_tiers
		WHERE product_id = $1::uuid AND min_quantity <= $2
		ORDER BY min_quantity DESC
		LIMIT 1`, productID, quantity).Scan(&tierK)
	if err == nil {
		return tierK, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	var base int64
	err = r.db.QueryRowContext(ctx, `SELECT price_kopeks FROM products WHERE id = $1::uuid`, productID).Scan(&base)
	return base, err
}

func (r *repository) GetProductImages(ctx context.Context, productID string) ([]models.ProductImage, error) {
	var imgs []models.ProductImage
	err := r.db.SelectContext(ctx, &imgs, `
		SELECT id, product_id, image_url, is_main, sort_order
		FROM product_images
		WHERE product_id = $1
		ORDER BY sort_order, is_main DESC`, productID)
	return imgs, err
}

func (r *repository) AddProductImage(ctx context.Context, productID string, imageURL string, isMain bool) (models.ProductImage, error) {
	var maxOrder int
	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(sort_order), -1) FROM product_images WHERE product_id = $1`, productID).Scan(&maxOrder)

	var img models.ProductImage
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO product_images (product_id, image_url, is_main, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, image_url, is_main, sort_order`,
		productID, imageURL, isMain, maxOrder+1).Scan(
		&img.ID, &img.ProductID, &img.ImageURL, &img.IsMain, &img.SortOrder,
	)
	return img, err
}

func (r *repository) RemoveProductImage(ctx context.Context, imageID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_images WHERE id = $1`, imageID)
	return err
}

func (r *repository) SetProductImages(ctx context.Context, productID string, images []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM product_images WHERE product_id = $1`, productID)
	if err != nil {
		return err
	}
	for i, url := range images {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO product_images (product_id, image_url, is_main, sort_order)
			VALUES ($1, $2, $3, $4)`, productID, url, i == 0, i)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repository) GetProductTags(ctx context.Context, productID string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.SelectContext(ctx, &tags, `
		SELECT t.id, t.name, t.color, t.created_at
		FROM tags t
		JOIN product_tags pt ON pt.tag_id = t.id
		WHERE pt.product_id = $1
		ORDER BY t.name`, productID)
	return tags, err
}

func (r *repository) AddProductTag(ctx context.Context, productID string, tagID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO product_tags (product_id, tag_id)
		VALUES ($1::uuid, $2::uuid)
		ON CONFLICT (product_id, tag_id) DO NOTHING`, productID, tagID)
	return err
}

func (r *repository) RemoveProductTag(ctx context.Context, productID string, tagID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM product_tags WHERE product_id = $1::uuid AND tag_id = $2::uuid`, productID, tagID)
	return err
}
