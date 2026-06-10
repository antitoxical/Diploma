package notifications

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, userUUID string, title, message string) (models.Notification, error)
	GetByUser(ctx context.Context, telegramID int64) ([]models.Notification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, telegramID int64) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, userUUID string, title, message string) (models.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, title, message, created_at)
		VALUES ($1, $2, $3, now())
		RETURNING id, user_id, title, message, is_read, created_at`

	var n models.Notification
	err := r.db.QueryRowContext(ctx, query, userUUID, title, message).Scan(
		&n.ID, &n.UserID, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt,
	)
	return n, err
}

func (r *repository) GetByUser(ctx context.Context, telegramID int64) ([]models.Notification, error) {
	query := `
		SELECT n.id, n.user_id, n.title, n.message, n.is_read, n.created_at
		FROM notifications n
		JOIN users u ON u.id = n.user_id
		WHERE u.telegram_id = $1
		ORDER BY n.created_at DESC`

	var list []models.Notification
	err := r.db.SelectContext(ctx, &list, query, telegramID)
	return list, err
}

func (r *repository) MarkRead(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read = true WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) MarkAllRead(ctx context.Context, telegramID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications n SET is_read = true
		FROM users u
		WHERE n.user_id = u.id AND u.telegram_id = $1 AND n.is_read = false`,
		telegramID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notifications WHERE id = $1`, id)
	return err
}
