package models

import "time"

type User struct {
	ID              string    `db:"id" json:"id"`
	TelegramID      int64     `db:"telegram_id" json:"telegram_id"`
	Username        *string   `db:"username" json:"username,omitempty"`
	FirstName       *string   `db:"first_name" json:"first_name,omitempty"`
	LastName        *string   `db:"last_name" json:"last_name,omitempty"`
	Phone           *string   `db:"phone" json:"phone,omitempty"`
	DeliveryAddress *string   `db:"delivery_address" json:"delivery_address,omitempty"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type CreateUser struct {
	TelegramID int64   `json:"telegram_id"`
	Username   *string `json:"username,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Phone      *string `json:"phone,omitempty"`
}

type UpdateUserProfile struct {
	Username        *string `json:"username,omitempty"`
	FirstName       *string `json:"first_name,omitempty"`
	LastName        *string `json:"last_name,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	DeliveryAddress *string `json:"delivery_address,omitempty"`
}

type UserAddress struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Label       *string   `db:"label" json:"label,omitempty"`
	FullAddress string    `db:"full_address" json:"full_address"`
	IsDefault   bool      `db:"is_default" json:"is_default"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type CreateUserAddress struct {
	TelegramID  int64  `json:"user_id"`
	Label       string `json:"label"`
	FullAddress string `json:"full_address"`
	IsDefault   bool   `json:"is_default"`
}
