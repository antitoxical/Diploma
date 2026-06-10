package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/internal/repository/users"
	"telegramshop_backend/pkg/logger"
	"telegramshop_backend/pkg/validate"
)

type Service interface {
	GetUserByTelegramID(ctx context.Context, telegramID int64) (models.User, error)
	CreateUser(ctx context.Context, input models.CreateUser) (models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, telegramID int64) error
	UpdateProfile(ctx context.Context, telegramID int64, input models.UpdateUserProfile) error

	CreateAddress(ctx context.Context, input models.CreateUserAddress) (models.UserAddress, error)
	GetAddresses(ctx context.Context, telegramID int64) ([]models.UserAddress, error)
	DeleteAddress(ctx context.Context, addressID string) error
}

type service struct {
	repo users.Repository
}

func NewService(repo users.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserByTelegramID(ctx context.Context, telegramID int64) (models.User, error) {
	user, err := s.repo.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, err
		}
		logger.Errorf("[GetUserByTelegramID] %v", err)
		return models.User{}, err
	}
	return user, nil
}

func (s *service) CreateUser(ctx context.Context, input models.CreateUser) (models.User, error) {
	return s.repo.CreateUser(ctx, input)
}

func (s *service) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) DeleteUser(ctx context.Context, telegramID int64) error {
	return s.repo.DeleteUser(ctx, telegramID)
}

func (s *service) UpdateProfile(ctx context.Context, telegramID int64, input models.UpdateUserProfile) error {
	if err := validate.UpdateUserProfile(input.FirstName, input.LastName, input.Phone); err != nil {
		return err
	}
	if input.DeliveryAddress != nil && strings.TrimSpace(*input.DeliveryAddress) != "" {
		if err := validate.Address(*input.DeliveryAddress, true); err != nil {
			return err
		}
	}
	return s.repo.UpdateUserProfile(ctx, telegramID, input)
}

func (s *service) CreateAddress(ctx context.Context, input models.CreateUserAddress) (models.UserAddress, error) {
	if err := validate.Address(input.FullAddress, true); err != nil {
		return models.UserAddress{}, err
	}
	addrs, err := s.repo.GetAddresses(ctx, input.TelegramID)
	if err != nil {
		return models.UserAddress{}, fmt.Errorf("address count: %w", err)
	}
	if len(addrs) >= validate.MaxAddresses {
		return models.UserAddress{}, fmt.Errorf("address limit reached (%d)", validate.MaxAddresses)
	}
	userUUID, err := s.repo.GetInternalIDByTelegramID(ctx, input.TelegramID)
	if err != nil {
		return models.UserAddress{}, fmt.Errorf("user lookup: %w", err)
	}
	return s.repo.CreateAddress(ctx, userUUID, input)
}

func (s *service) GetAddresses(ctx context.Context, telegramID int64) ([]models.UserAddress, error) {
	return s.repo.GetAddresses(ctx, telegramID)
}

func (s *service) DeleteAddress(ctx context.Context, addressID string) error {
	return s.repo.DeleteAddress(ctx, addressID)
}
