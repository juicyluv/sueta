package db

import (
	"context"

	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// db implementes user storage interface.
type db struct {
	logger     logger.Logger
	collection *mongo.Collection
}

// New Storage returns a new user storage instance.
func NewStorage(storage *mongo.Database, collection string) user.Storage {
	return &db{
		logger:     logger.GetLogger(),
		collection: storage.Collection(collection),
	}
}

// Create inserts a new row in the database.
// Returns an error on failure and token string on success.
func (d *db) Create(ctx context.Context, user *user.CreateUserDTO) (string, error) {
	return "", nil
}

// FindByEmail finds the user by given email.
// Returns an error on failure or No Rows Error if there's no user with given email.
func (d *db) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

// FindById finds the user by given uuid.
// Returns an error on failure or No Rows Error if there's no user with given uuid.
func (d *db) FindById(ctx context.Context, uuid string) (*user.User, error) {
	return nil, nil
}

// UpdatePartially updates the user with given uuid.
// Returns an error if something went wrong or No Rows error if
// there's no user with given uuid.
func (d *db) UpdatePartially(ctx context.Context, user *user.UpdateUserDTO) error {
	return nil
}

// Delete deletes the user row with given uuid.
// Returns an error on failure. Returns nil if there's no such user with given uuid.
func (d *db) Delete(ctx context.Context, uuid string) error {
	return nil
}
