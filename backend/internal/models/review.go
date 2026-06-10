package models

import "time"

type Review struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	Rating      int       `db:"rating" json:"rating"`
	Comment     *string   `db:"comment" json:"comment,omitempty"`
	IsModerated bool      `db:"is_moderated" json:"is_moderated"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type UpsertReviewInput struct {
	TelegramID int64  `json:"user_id"`
	ProductID  string `json:"product_id"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
}
