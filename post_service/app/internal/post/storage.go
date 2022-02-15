package post

import "context"

// Storage descibes a post storage functionality.
type Storage interface {
	Create(ctx context.Context, post *Post) (string, error)
	FindById(ctx context.Context, uuid string) (*Post, error)
	UpdatePartially(ctx context.Context, post *Post) error
	Delete(ctx context.Context, uuid string) error
}
