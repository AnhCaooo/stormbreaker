// AnhCao 2024
package db

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

// Function to connect to mongo database instance and create collection if it does not exist
func Init(ctx context.Context, cfg models.Database) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(getURI(cfg))
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

func getURI(cfg models.Database) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/?timeoutMS=5000", cfg.Name, cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}
