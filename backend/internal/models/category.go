package models

import "time"

type Category struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ParentID  *string   `db:"parent_id" json:"parent_id,omitempty"`
	ImageURL  *string   `db:"image_url" json:"image_url,omitempty"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UpdateCategoryInput struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
