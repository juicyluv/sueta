package user

import "context"

type Storage interface {
	Create(ctx context.Context, user *CreateUserDTO) (string, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, uuid string) (*User, error)
	UpdatePartially(ctx context.Context, user *UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}
