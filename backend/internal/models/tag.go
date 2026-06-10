package models

import "time"

type Tag struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Color     *string   `db:"color" json:"color,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateTagInput struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

type ProductTag struct {
	ID        string `db:"id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	TagID     string `db:"tag_id" json:"tag_id"`
}
