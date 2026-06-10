package models

import "time"

type SearchHistoryEntry struct {
	ID           string    `db:"id" json:"id"`
	UserID       *string   `db:"user_id" json:"user_id,omitempty"`
	Query        string    `db:"query" json:"query"`
	ResultsCount int       `db:"results_count" json:"results_count"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
