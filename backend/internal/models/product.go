package models

import "time"

type Product struct {
	ID            string         `db:"id" json:"id"`
	CategoryID    string         `db:"category_id" json:"category_id"`
	FirmID        *string        `db:"firm_id" json:"firm_id,omitempty"`
	Name          string         `db:"name" json:"name"`
	Description   *string        `db:"description" json:"description,omitempty"`
	PriceKopeks   int64          `db:"price_kopeks" json:"price_kopeks"`
	StockQuantity int            `db:"stock_quantity" json:"stock_quantity"`
	IsActive      bool           `db:"is_active" json:"is_active"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	Images        []ProductImage `db:"-" json:"images,omitempty"`
	Tags          []Tag          `db:"-" json:"tags,omitempty"`
	// Список URL из админки (POST /admin/products); не колонка БД.
	ImageURLs []string `json:"image_urls,omitempty" db:"-"`
}

type ProductImage struct {
	ID        string `db:"id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	ImageURL  string `db:"image_url" json:"image_url"`
	IsMain    bool   `db:"is_main" json:"is_main"`
	SortOrder int    `db:"sort_order" json:"sort_order"`
}

type UpdateProductInput struct {
	Name          string  `json:"name"`
	CategoryID    string  `json:"category_id"`
	FirmID        *string `json:"firm_id,omitempty"`
	Description   *string `json:"description,omitempty"`
	PriceKopeks   int64   `json:"price_kopeks"`
	StockQuantity int     `json:"stock_quantity"`
	IsActive      *bool   `json:"is_active,omitempty"`
	// nil — не трогать картинки; не-nil (в т.ч. []) — заменить набор целиком.
	ImageURLs *[]string `json:"image_urls,omitempty"`
}

type ProductFilter struct {
	CategoryID           *string  `query:"category_id"`
	IncludeSubcategories bool     `query:"include_subcategories"`
	FirmID               *string  `query:"firm_id"`
	TagID                *string  `query:"tag_id"`
	TagIDs               []string `query:"-"`
	Search               *string  `query:"search"`
	MinPrice             *int64   `query:"min_price"`
	MaxPrice             *int64   `query:"max_price"`
	InStock              *bool    `query:"in_stock"`
	MinRating            *float64 `query:"min_rating"`
	NewDays              *int     `query:"new_days"`
	HasReviews           *bool    `query:"has_reviews"`
	SortBy               string   `query:"sort_by"`
	Page                 int      `query:"page"`
	Limit                int      `query:"limit"`
}

type PaginatedProducts struct {
	Products   []Product `json:"products"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

type ImageInput struct {
	Image string `json:"image" example:"https://example.com/image.jpg"`
}

type ImagesInput struct {
	Images []string `json:"images"`
}

type StockInput struct {
	Stock int `json:"stock"`
}
