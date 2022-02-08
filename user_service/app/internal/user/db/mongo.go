package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/juicyluv/sueta/user_service/app/internal/user/apperror"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Check whether db implements user storage interface.
var _ user.Storage = &db{}

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
// It uses InsertOne behind the scenes.
// Returns an error on failure or inserted user uuid on success.
func (d *db) Create(ctx context.Context, user *user.User) (string, error) {
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		e := fmt.Errorf("cannot insert user in database: %w", err)
		d.logger.Warn(e)
		return "", e
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		e := fmt.Errorf("cannot convert user uuid to object id: %w", err)
		d.logger.Warn(e)
		return "", e
	}

	return id.Hex(), nil
}

// FindByEmail finds the user by given email.
// Returns user instance on success, but on failure
// returns an error or No Rows Error
// if there's no user with given email.
func (d *db) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}

// FindById finds the user by given uuid.
// Returns user instance on success, but on failure
// returns an error or No Rows Error if there's no user with given uuid.
func (d *db) FindById(ctx context.Context, uuid string) (*user.User, error) {
	var user *user.User

	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hex to objectid: %w", err)
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, apperror.ErrNotFound
		}
		d.logger.Error(result.Err())
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if err = result.Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	return user, nil
}

// UpdatePartially updates the user with given uuid.
// Returns an error if something went wrong or No Rows error if
// there's no user with given uuid.
func (d *db) UpdatePartially(ctx context.Context, user *user.User) error {
	return nil
}

// Delete deletes the user row with given uuid.
// Returns an error on failure. Returns nil if there's no such user with given uuid.
func (d *db) Delete(ctx context.Context, uuid string) error {
	return nil
}
