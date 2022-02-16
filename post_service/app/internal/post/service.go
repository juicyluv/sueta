package post

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicyluv/sueta/post_service/app/internal/post/apperror"
	"github.com/juicyluv/sueta/post_service/app/pkg/logger"
)

// Service describes post service functionality.
type Service interface {
	Create(ctx context.Context, post *CreatePostDTO) (string, error)
	GetById(ctx context.Context, uuid string) (*Post, error)
	UpdatePartially(ctx context.Context, user *UpdatePostDTO) error
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
func (s *service) Create(ctx context.Context, input *CreatePostDTO) (string, error) {
	post := &Post{
		Title:     input.Title,
		Content:   input.Content,
		UserUUID:  input.UserUUID,
		CreatedAt: time.Now().UTC().Format("2006/01/02"),
		UpdatedAt: time.Now().UTC().Format("2006/01/02"),
	}

	id, err := s.storage.Create(ctx, post)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetById will find a user with specified uuid in storage.
// Returns an error on failure of there's no user with this uuid.
func (s *service) GetById(ctx context.Context, uuid string) (*Post, error) {
	post, err := s.storage.FindById(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNoRows) || errors.Is(err, apperror.ErrInvalidUUID) {
			return nil, err
		}
		err = fmt.Errorf("failed to find post by uuid: %v", err)
		s.logger.Warn(err)
		return nil, err
	}

	return post, nil
}

// UpdatePartially will find the user with provided uuid.
// If there is no user with such id, returns No Rows error.
// Then passwords will be compared. If it don't match, returns
// Wrong Password error. Then updates the user. If something went wrong,
// returns an error and nil if everything is OK.
func (s *service) UpdatePartially(ctx context.Context, post *UpdatePostDTO) error {
	p, err := s.GetById(ctx, post.UUID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNoRows) {
			s.logger.Warnf("failed to get the post: %v", err)
		}
		return err
	}

	if post.Title != nil {
		p.Title = *post.Title
	}

	if post.Content != nil {
		p.Content = *post.Content
	}

	if post.UserUUID != nil {
		p.UserUUID = *post.UserUUID
	}

	err = s.storage.UpdatePartially(ctx, p)
	if err != nil {
		s.logger.Warnf("failed to update the post: %v", err)
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
			s.logger.Warnf("failed to delete the post: %v", err)
		}
		return err
	}

	return nil
}
