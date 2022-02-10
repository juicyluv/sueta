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

// Create will check whether provided email already taken.
// If it is, returns an error. Then it will hash user password
// and try to insert the user. Returns inserted UUID or an error
// on failure.
func (s *service) Create(ctx context.Context, input *CreateUserDTO) (string, error) {
	found, err := s.storage.FindByEmail(ctx, input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			return "", err
		}
	}

	if found != nil {
		return "", apperror.ErrEmailTaken
	}

	user := &User{
		Email:        input.Email,
		Username:     input.Username,
		Password:     input.Password,
		Verified:     false,
		RegisteredAt: time.Now().UTC().Format("2006/01/02"),
	}

	if err := user.HashPassword(); err != nil {
		s.logger.Warn("could not encrypt user password: %v", err)
		return "", err
	}

	id, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetByEmailAndPassword will find a user with provided email.
// If there's no such user with this email, returns No Rows error.
// If password doesn't match, returns Wrong Password error.
// Returns a user if everything is OK.
func (s *service) GetByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	user, err := s.storage.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		s.logger.Warn("error occurred on finding user by email: %v", err)
		return nil, err
	}

	if !user.ComparePassword(password) {
		return nil, apperror.ErrWrongPassword
	}

	return user, nil
}

// GetById will find a user with specified uuid in storage.
// Returns an error on failure of there's no user with this uuid.
func (s *service) GetById(ctx context.Context, uuid string) (*User, error) {
	user, err := s.storage.FindById(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) {
			return nil, err
		}
		err = fmt.Errorf("failed to find user by uuid: %v", err)
		s.logger.Warn(err)
		return nil, err
	}

	return user, nil
}

// UpdatePartially will find the user with provided uuid.
// If there is no user with such id, returns No Rows error.
// Then passwords will be compared. If it don't match, returns
// Wrong Password error. Then updates the user. If something went wrong,
// returns an error and nil if everything is OK.
func (s *service) UpdatePartially(ctx context.Context, user *UpdateUserDTO) error {
	u, err := s.GetById(ctx, user.UUID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warn("failed to get the user: %v", err)
		}
		return err
	}

	if !u.ComparePassword(*user.OldPassword) {
		return apperror.ErrWrongPassword
	}

	if user.NewPassword != nil {
		u.Password = *user.NewPassword
		err = u.HashPassword()
		if err != nil {
			s.logger.Warn("failed ot hash password: %v", err)
			return err
		}
	}

	if user.Email != nil {
		u.Email = *user.Email
	}

	if user.Username != nil {
		u.Username = *user.Username
	}

	if err := user.Validate(); err != nil {
		return err
	}

	err = s.storage.UpdatePartially(ctx, u)
	if err != nil {
		s.logger.Warn("failed to update the user: %v", err)
		return err
	}

	return nil
}

// Delete tries to delete the user with provided uuid.
// Returns an error on failure or nil if query has been executed.
func (s *service) Delete(ctx context.Context, uuid string) error {
	err := s.storage.Delete(ctx, uuid)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warn("failed to delete the user: %v", err)
		}
	}

	return nil
}
