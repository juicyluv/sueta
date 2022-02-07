package db

import (
	"context"

	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	logger     logger.Logger
	collection *mongo.Collection
}

func NewStorage(storage *mongo.Database, collection string) user.Storage {
	return &db{
		logger:     logger.GetLogger(),
		collection: storage.Collection(collection),
	}
}

func (d *db) Create(ctx context.Context, user *user.CreateUserDTO) (string, error) {
	return "", nil
}

func (d *db) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

func (d *db) FindById(ctx context.Context, uuid string) (*user.User, error) {
	return nil, nil
}

func (d *db) UpdatePartially(ctx context.Context, user *user.UpdateUserDTO) error {
	return nil
}

func (d *db) Delete(ctx context.Context, uuid string) error {
	return nil
}
