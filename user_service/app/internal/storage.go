package internal

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, uuid string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, uuid string) error
}
