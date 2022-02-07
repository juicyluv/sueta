package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient tries to connect to database with provided options.
// If everything is OK, returns a new mongo client instance.
// Returns an error if something went wrong.
func NewMongoClient(ctx context.Context, host, port, username, password, database, authDB string) (*mongo.Database, error) {
	var mongoURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	clientOptions := options.Client().ApplyURI(mongoURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}

		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	}

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
