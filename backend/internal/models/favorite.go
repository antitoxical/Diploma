package models

import "time"

type Favorite struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateFavorite struct {
	UserID    int64  `json:"user_id"` // telegram_id
	ProductID string `json:"product_id"`
}

type DeleteFavorite struct {
	UserID    int64  `json:"user_id"`
	ProductID string `json:"product_id"`
}
