package models

import "time"

type Notification struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Title     *string   `db:"title" json:"title,omitempty"`
	Message   *string   `db:"message" json:"message,omitempty"`
	IsRead    bool      `db:"is_read" json:"is_read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateNotification struct {
	UserID  int64  `json:"user_id"` // telegram_id
	Title   string `json:"title"`
	Message string `json:"message"`
}
