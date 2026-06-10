package models

import "time"

type Firm struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UpdateFirmInput struct {
	Name string `db:"name" json:"name"`
}
