package admins

import (
	"context"
	"errors"
	"time"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/admins"
	"telegramshop_backend/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

var ErrInvalidCredentials = errors.New("invalid username or password")

type Service interface {
	Login(ctx context.Context, input models.AdminLoginInput, jwtSecret string) (models.AdminTokenResponse, error)
	GetAll(ctx context.Context) ([]models.Admin, error)
	GetByID(ctx context.Context, id string) (models.Admin, error)
	CreateAdmin(ctx context.Context, input models.AdminCreateInput) (models.Admin, error)
	UpdateAdmin(ctx context.Context, id string, input models.AdminUpdateInput) (models.Admin, error)
	GetAuditLogs(ctx context.Context, limit int) ([]models.AuditLog, error)
	WriteAuditLog(ctx context.Context, adminID, action, tableName string, recordID *string, details map[string]any) error
}

type service struct {
	repo admins.Repository
}

func NewService(repo admins.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Login(ctx context.Context, input models.AdminLoginInput, jwtSecret string) (models.AdminTokenResponse, error) {
	admin, err := s.repo.GetByUsername(ctx, input.Username)
	if err != nil {
		return models.AdminTokenResponse{}, ErrInvalidCredentials
	}

	if !admin.IsActive {
		return models.AdminTokenResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		return models.AdminTokenResponse{}, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(admin.ID, admin.Username, jwtSecret, tokenTTL)
	if err != nil {
		return models.AdminTokenResponse{}, err
	}

	return models.AdminTokenResponse{
		Token: token,
		Admin: admin,
	}, nil
}

func (s *service) GetAll(ctx context.Context) ([]models.Admin, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (models.Admin, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateAdmin(ctx context.Context, id string, input models.AdminUpdateInput) (models.Admin, error) {
	var hash string
	if input.Password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.Admin{}, err
		}
		hash = string(h)
	}
	return s.repo.Update(ctx, id, input.Username, hash, input.IsActive)
}

func (s *service) CreateAdmin(ctx context.Context, input models.AdminCreateInput) (models.Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Admin{}, err
	}
	return s.repo.Create(ctx, input.Username, string(hash))
}

func (s *service) GetAuditLogs(ctx context.Context, limit int) ([]models.AuditLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	return s.repo.GetAuditLogs(ctx, limit)
}

func (s *service) WriteAuditLog(ctx context.Context, adminID, action, tableName string, recordID *string, details map[string]any) error {
	return s.repo.WriteAuditLog(ctx, adminID, action, tableName, recordID, details)
}
