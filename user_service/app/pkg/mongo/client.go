package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient tries to connect to database with provided options.
// If everything is OK, returns a new mongo database instance.
// Returns an error if something went wrong.
func NewMongoClient(ctx context.Context, database, mongoURL string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongodb: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not ping mongodb: %w", err)
	}

	return client.Database(database), nil
}
