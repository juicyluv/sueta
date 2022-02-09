package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicyluv/sueta/user_service/app/internal/user/apperror"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
)

// Service describes user service functionality.
type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (string, error)
	GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
	GetById(ctx context.Context, uuid string) (*User, error)
	UpdatePartially(ctx context.Context, user *UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

type service struct {
	logger  logger.Logger
	storage Storage
}

// NewService returns a new instance that implements Service interface.
func NewService(storage Storage, logger logger.Logger) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

func (s *service) Create(ctx context.Context, input *CreateUserDTO) (string, error) {
	_, err := s.storage.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, apperror.ErrNoRows) {
		return "", err
	}

	user := &User{
		Email:        input.Email,
		Username:     input.Username,
		Password:     input.Password,
		Verified:     false,
		RegisteredAt: time.Now().UTC().String(),
	}

	if err := user.HashPassword(); err != nil {
		s.logger.Warn("could not encrypt user password: %w", err)
		return "", err
	}

	id, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	user, err := s.storage.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warn("error occurred on finding user by email: %w", err)
		return nil, err
	}

	if !user.ComparePassword(password) {
		return nil, errors.New("password doesn't match")
	}

	return user, nil
}

// GetById will find a user with specified uuid in storage.
// Returns an error on failure of there's no user with this uuid.
func (s *service) GetById(ctx context.Context, uuid string) (*User, error) {
	user, err := s.storage.FindById(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return user, err
		}
		err = fmt.Errorf("failed to find user by uuid: %v", err)
		s.logger.Warn(err)
		return nil, err
	}

	return user, nil
}

func (s *service) UpdatePartially(ctx context.Context, user *UpdateUserDTO) error {
	return nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	return nil
}
