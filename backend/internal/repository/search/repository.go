package search

import (
	"context"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Record(ctx context.Context, userUUID *string, query string, resultsCount int) error
	GetUserHistory(ctx context.Context, telegramID int64, limit int) ([]models.SearchHistoryEntry, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Record(ctx context.Context, userUUID *string, query string, resultsCount int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO search_history (user_id, query, results_count, created_at)
		VALUES ($1, $2, $3, now())`, userUUID, query, resultsCount)
	return err
}

func (r *repository) GetUserHistory(ctx context.Context, telegramID int64, limit int) ([]models.SearchHistoryEntry, error) {
	var entries []models.SearchHistoryEntry
	err := r.db.SelectContext(ctx, &entries, `
		SELECT sh.id, sh.user_id, sh.query, sh.results_count, sh.created_at
		FROM search_history sh
		JOIN users u ON u.id = sh.user_id
		WHERE u.telegram_id = $1
		ORDER BY sh.created_at DESC
		LIMIT $2`, telegramID, limit)
	return entries, err
}
