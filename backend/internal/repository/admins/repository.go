package admins

import (
	"context"
	"database/sql"
	"encoding/json"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (models.Admin, error)
	GetByID(ctx context.Context, id string) (models.Admin, error)
	Create(ctx context.Context, username, passwordHash string) (models.Admin, error)
	Update(ctx context.Context, id, username, passwordHash string, isActive bool) (models.Admin, error)
	GetAll(ctx context.Context) ([]models.Admin, error)

	CreateSession(ctx context.Context, adminID, tokenHash, ipAddress, userAgent string, expiresAt string) error
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteExpiredSessions(ctx context.Context) error

	WriteAuditLog(ctx context.Context, adminID, action, tableName string, recordID *string, details map[string]any) error
	GetAuditLogs(ctx context.Context, limit int) ([]models.AuditLog, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByUsername(ctx context.Context, username string) (models.Admin, error) {
	query := `SELECT id, username, password_hash, is_active, created_at FROM admins WHERE username = $1`
	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, username).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.IsActive, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Admin{}, sql.ErrNoRows
	}
	return a, err
}

func (r *repository) GetByID(ctx context.Context, id string) (models.Admin, error) {
	query := `SELECT id, username, password_hash, is_active, created_at FROM admins WHERE id = $1`
	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.IsActive, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Admin{}, sql.ErrNoRows
	}
	return a, err
}

func (r *repository) Update(ctx context.Context, id, username, passwordHash string, isActive bool) (models.Admin, error) {
	var query string
	var args []interface{}
	if passwordHash != "" {
		query = `UPDATE admins SET username=$2, password_hash=$3, is_active=$4 WHERE id=$1
			RETURNING id, username, password_hash, is_active, created_at`
		args = []interface{}{id, username, passwordHash, isActive}
	} else {
		query = `UPDATE admins SET username=$2, is_active=$3 WHERE id=$1
			RETURNING id, username, password_hash, is_active, created_at`
		args = []interface{}{id, username, isActive}
	}
	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.IsActive, &a.CreatedAt)
	return a, err
}

func (r *repository) Create(ctx context.Context, username, passwordHash string) (models.Admin, error) {
	query := `
		INSERT INTO admins (username, password_hash, created_at)
		VALUES ($1, $2, now())
		RETURNING id, username, password_hash, is_active, created_at`
	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, username, passwordHash).Scan(
		&a.ID, &a.Username, &a.PasswordHash, &a.IsActive, &a.CreatedAt,
	)
	return a, err
}

func (r *repository) GetAll(ctx context.Context) ([]models.Admin, error) {
	var admins []models.Admin
	err := r.db.SelectContext(ctx, &admins, `
		SELECT id, username, password_hash, is_active, created_at
		FROM admins ORDER BY created_at DESC`)
	return admins, err
}

// ── Sessions ────────────────────────────────────────────────────────

func (r *repository) CreateSession(ctx context.Context, adminID, tokenHash, ipAddress, userAgent string, expiresAt string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO admin_sessions (admin_id, token_hash, ip_address, user_agent, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5::timestamptz, now())`,
		adminID, tokenHash, ipAddress, userAgent, expiresAt)
	return err
}

func (r *repository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM admin_sessions WHERE id = $1`, sessionID)
	return err
}

func (r *repository) DeleteExpiredSessions(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM admin_sessions WHERE expires_at < now()`)
	return err
}

// ── Audit Log ───────────────────────────────────────────────────────

func (r *repository) WriteAuditLog(ctx context.Context, adminID, action, tableName string, recordID *string, details map[string]any) error {
	var detailsJSON any
	if len(details) > 0 {
		b, err := json.Marshal(details)
		if err != nil {
			return err
		}
		detailsJSON = b
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO audit_logs (admin_id, action, table_name, record_id, details, created_at)
		VALUES ($1, $2, $3, $4, $5, now())`,
		adminID, action, tableName, recordID, detailsJSON)
	return err
}

func (r *repository) GetAuditLogs(ctx context.Context, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.SelectContext(ctx, &logs, `
		SELECT al.id, al.admin_id, a.username AS admin_username,
		       al.action, al.table_name, al.record_id, al.details, al.created_at
		FROM audit_logs al
		LEFT JOIN admins a ON a.id = al.admin_id
		ORDER BY al.created_at DESC
		LIMIT $1`, limit)
	return logs, err
}
