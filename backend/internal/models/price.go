package models

import "time"

// ProductPriceTier represents a volume-based price tier (from min_quantity units).
type ProductPriceTier struct {
	ID          string    `db:"id" json:"id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	MinQuantity int       `db:"min_quantity" json:"min_quantity"`
	PriceKopeks int64     `db:"price_kopeks" json:"price_kopeks"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type UpdatePriceTierInput struct {
	MinQuantity int   `json:"min_quantity"`
	PriceKopeks int64 `json:"price_kopeks"`
}

type UpdateTierMinQuantity struct {
	NewMinQuantity int `json:"new_min_quantity" example:"5"`
}
