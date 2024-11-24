// AnhCao 2024
package db

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var Collection *mongo.Collection

// Init to connect to mongo database instance and create collection if it does not exist
func Init(ctx context.Context, cfg models.Database) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(getURI(cfg))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err.Error())
	}

	Collection = client.Database(cfg.Name).Collection(cfg.Collection)

	logger.Logger.Info("Successfully connected to database")
	return client, nil
}

func getURI(cfg models.Database) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?timeoutMS=5000", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}

// InsertPriceSettings inserts a new document into the PriceSettings collection
//
// Use case: set the price settings after user signed up.
func InsertPriceSettings(ctx context.Context, settings models.PriceSettings) error {
	_, err := Collection.InsertOne(ctx, settings)
	return err
}

// GetPriceSettings retrieves a document by UserID
func GetPriceSettings(ctx context.Context, userID string) (*models.PriceSettings, error) {
	filter := bson.M{"user_id": userID}

	var settings models.PriceSettings
	err := Collection.FindOne(ctx, filter).Decode(&settings)
	if err != nil {
		return nil, fmt.Errorf("failed to get price setting: %s", err.Error())
	}
	logger.Logger.Info("get price settings successfully", zap.Any("user_id", userID))
	return &settings, nil
}

// PatchPriceSettings updates partial data for user's price settings.
//
// Use case: update the price settings for specific user
func PatchPriceSettings(ctx context.Context, settings models.PriceSettings) error {
	filter := bson.M{"user_id": settings.UserID}

	updates := bson.M{
		"$set": bson.M{
			"vat_included": settings.VatIncluded,
			"margin":       settings.Marginal,
		},
	}
	result, err := Collection.UpdateOne(ctx, filter, updates)
	if err != nil {
		return fmt.Errorf("failed to update price settings: %s", err.Error())
	}
	logger.Logger.Info("update price settings successfully", zap.Any("updated_amount", result.MatchedCount))
	return nil
}

// DeletePriceSettings deletes user's price settings.
//
// Use case: delete the price settings if user is deleted
func DeletePriceSettings(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}

	result, err := Collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete price settings. Error: %s", err.Error())
	}
	logger.Logger.Info("delete user price settings successfully", zap.Int64("deleted_amount", result.DeletedCount))
	return nil
}
