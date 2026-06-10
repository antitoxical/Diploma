package users

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.CreateUser) (models.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (models.User, error)
	GetInternalIDByTelegramID(ctx context.Context, telegramID int64) (string, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	UpdateUserProfile(ctx context.Context, telegramID int64, u models.UpdateUserProfile) error
	DeleteUser(ctx context.Context, telegramID int64) error
	GetAll(ctx context.Context) ([]models.User, error)
	ListTelegramIDs(ctx context.Context) ([]int64, error)

	CreateAddress(ctx context.Context, userUUID string, input models.CreateUserAddress) (models.UserAddress, error)
	GetAddresses(ctx context.Context, telegramID int64) ([]models.UserAddress, error)
	DeleteAddress(ctx context.Context, addressID string) error
}

const userCols = `id, telegram_id, username, first_name, last_name, phone, delivery_address, is_active, created_at`

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user models.CreateUser) (models.User, error) {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, phone, created_at)
		VALUES ($1, $2, $3, $4, $5, now())
		RETURNING ` + userCols

	var out models.User
	err := r.db.QueryRowContext(ctx, query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Phone,
	).Scan(
		&out.ID,
		&out.TelegramID,
		&out.Username,
		&out.FirstName,
		&out.LastName,
		&out.Phone,
		&out.DeliveryAddress,
		&out.IsActive,
		&out.CreatedAt,
	)
	return out, err
}

func (r *repository) GetUserByTelegramID(ctx context.Context, telegramID int64) (models.User, error) {
	query := `SELECT ` + userCols + ` FROM users WHERE telegram_id = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.DeliveryAddress,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *repository) GetInternalIDByTelegramID(ctx context.Context, telegramID int64) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE telegram_id = $1`, telegramID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT ` + userCols + ` FROM users WHERE username = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.DeliveryAddress,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.User{}, nil
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *repository) UpdateUserProfile(ctx context.Context, telegramID int64, u models.UpdateUserProfile) error {
	query := `
		UPDATE users
		SET username = COALESCE($1, username),
		    first_name = COALESCE($2, first_name),
		    last_name = COALESCE($3, last_name),
		    phone = COALESCE($4, phone),
		    delivery_address = COALESCE($5, delivery_address)
		WHERE telegram_id = $6`
	_, err := r.db.ExecContext(ctx, query, u.Username, u.FirstName, u.LastName, u.Phone, u.DeliveryAddress, telegramID)
	return err
}

func (r *repository) DeleteUser(ctx context.Context, telegramID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE telegram_id = $1`, telegramID)
	return err
}

func (r *repository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT ` + userCols + ` FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.DeliveryAddress,
			&user.IsActive,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *repository) ListTelegramIDs(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids, `
		SELECT telegram_id FROM users
		WHERE is_active = true AND telegram_id IS NOT NULL AND telegram_id != 0`)
	return ids, err
}

// ── Addresses ───────────────────────────────────────────────────────

func (r *repository) CreateAddress(ctx context.Context, userUUID string, input models.CreateUserAddress) (models.UserAddress, error) {
	if input.IsDefault {
		_, _ = r.db.ExecContext(ctx, `UPDATE user_addresses SET is_default = false WHERE user_id = $1`, userUUID)
	}
	var addr models.UserAddress
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO user_addresses (user_id, label, full_address, is_default, created_at)
		VALUES ($1, $2, $3, $4, now())
		RETURNING id, user_id, label, full_address, is_default, created_at`,
		userUUID, input.Label, input.FullAddress, input.IsDefault).Scan(
		&addr.ID, &addr.UserID, &addr.Label, &addr.FullAddress, &addr.IsDefault, &addr.CreatedAt,
	)
	return addr, err
}

func (r *repository) GetAddresses(ctx context.Context, telegramID int64) ([]models.UserAddress, error) {
	var addrs []models.UserAddress
	err := r.db.SelectContext(ctx, &addrs, `
		SELECT a.id, a.user_id, a.label, a.full_address, a.is_default, a.created_at
		FROM user_addresses a
		JOIN users u ON u.id = a.user_id
		WHERE u.telegram_id = $1
		ORDER BY a.is_default DESC, a.created_at DESC`, telegramID)
	return addrs, err
}

func (r *repository) DeleteAddress(ctx context.Context, addressID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_addresses WHERE id = $1`, addressID)
	return err
}
