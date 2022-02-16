package db

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/juicyluv/sueta/post_service/app/internal/post"
	"github.com/juicyluv/sueta/post_service/app/pkg/logger"
)

// Check whether db implements user storage interface.
var _ post.Storage = &db{}

// db implementes user storage interface.
type db struct {
	logger     logger.Logger
	connection *pgx.Conn
}

// New Storage returns a new user storage instance.
func NewStorage(storage *pgx.Conn) post.Storage {
	return &db{
		logger:     logger.GetLogger(),
		connection: storage,
	}
}

// Create inserts a new row in the database.
// It uses InsertOne behind the scenes.
// Returns an error on failure or inserted user uuid on success.
func (d *db) Create(ctx context.Context, post *post.Post) (string, error) {
	// result, err := d.collection.InsertOne(ctx, user)
	// if err != nil {
	// 	e := fmt.Errorf("cannot insert user in database: %w", err)
	// 	d.logger.Warn(e)
	// 	return "", e
	// }

	return "", nil
}

// FindById finds the user by given uuid.
// Returns user instance on success, but on failure
// returns an error or No Rows Error if there's no user with given uuid.
func (d *db) FindById(ctx context.Context, uuid string) (*post.Post, error) {
	// var post *post.Post

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	// result := d.collection.FindOne(ctx, filter)
	// if result.Err() != nil {
	// 	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
	// 		return nil, apperror.ErrNoRows
	// 	}
	// 	d.logger.Error(result.Err())
	// 	return nil, fmt.Errorf("failed to execute query: %w", err)
	// }

	return nil, nil
}

// UpdatePartially updates the user with new provided values.
// Returns an error if something went wrong or No Rows error if
// there's no user with given uuid.
func (d *db) UpdatePartially(ctx context.Context, post *post.Post) error {
	// result, err := d.collection.UpdateOne(ctx, filter, query)
	// if err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return apperror.ErrNoRows
	// 	}
	// 	d.logger.Warn("failed to execute query: %w", err)
	// 	return err
	// }

	return nil
}

// Delete deletes the user row with given uuid. Returns an error on failure.
// Returns ErrNoRows if there's no such user with given uuid.
func (d *db) Delete(ctx context.Context, uuid string) error {
	// result := d.collection.FindOneAndDelete(ctx, filter)
	// if err := result.Err(); err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return apperror.ErrNoRows
	// 	}
	// 	return fmt.Errorf("cannot delete user: %v", err)
	// }

	return nil
}
