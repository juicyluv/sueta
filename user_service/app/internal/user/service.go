package user

import (
	"context"

	"github.com/juicyluv/sueta/user_service/app/internal/logger"
)

type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (string, error)
	GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
	GetById(ctx context.Context, uuid string) (*User, error)
	UpdatePartially(ctx context.Context, user *UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

type service struct {
	logger  logger.Logging
	storage Storage
}

func NewService(storage Storage, logger logger.Logging) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

func (s *service) Create(ctx context.Context, user *CreateUserDTO) (string, error) {
	return "", nil
}

func (s *service) GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	return nil, nil
}

func (s *service) GetById(ctx context.Context, uuid string) (*User, error) {
	return nil, nil
}

func (s *service) UpdatePartially(ctx context.Context, user *UpdateUserDTO) error {
	return nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	return nil
}
