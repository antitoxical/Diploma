package models

import "time"

type BasketItem struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CreateBasketItem holds input data; user_id in JSON corresponds to telegram_id.
type CreateBasketItem struct {
	UserID    int64  `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type DeleteBasketItem struct {
	UserID    int64  `json:"user_id"`
	ProductID string `json:"product_id"`
}
