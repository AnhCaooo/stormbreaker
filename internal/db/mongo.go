package db

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function to connect to mongo database instance
func Init(ctx context.Context, URI string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database. Error: %s", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database. Error: %s", err.Error())
	}

	logger.Logger.Info("Successfully connected to database")
	return client, nil
}
