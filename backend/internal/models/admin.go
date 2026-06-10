package models

import (
	"encoding/json"
	"time"
)

type Admin struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type AdminLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminCreateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminUpdateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

type AdminTokenResponse struct {
	Token string `json:"token"`
	Admin Admin  `json:"admin"`
}

type AdminSession struct {
	ID        string    `db:"id" json:"id"`
	AdminID   string    `db:"admin_id" json:"admin_id"`
	TokenHash string    `db:"token_hash" json:"-"`
	IPAddress *string   `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent *string   `db:"user_agent" json:"user_agent,omitempty"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type AuditLog struct {
	ID            string           `db:"id" json:"id"`
	AdminID       *string          `db:"admin_id" json:"admin_id,omitempty"`
	AdminUsername *string          `db:"admin_username" json:"admin_username,omitempty"`
	Action        string           `db:"action" json:"action"`
	TableName     *string          `db:"table_name" json:"table_name,omitempty"`
	RecordID      *string          `db:"record_id" json:"record_id,omitempty"`
	Details       *json.RawMessage `db:"details" json:"details,omitempty"`
	CreatedAt     time.Time        `db:"created_at" json:"created_at"`
}
