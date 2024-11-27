// AnhCao 2024
package db

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Mongo struct {
	config           *models.Database
	logger           *zap.Logger
	ctx              context.Context
	collectionClient *mongo.Collection
}

func Init(ctx context.Context, config *models.Database, logger *zap.Logger, collectionClient *mongo.Collection) *Mongo {
	return &Mongo{
		config:           config,
		logger:           logger,
		ctx:              ctx,
		collectionClient: nil,
	}
}

// Init to connect to mongo database instance and create collection if it does not exist
func (db Mongo) EstablishConnection() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db.getURI())
	client, err := mongo.Connect(db.ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	err = client.Ping(db.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err.Error())
	}

	db.collectionClient = client.Database(db.config.Name).Collection(db.config.Collection)

	db.logger.Info("Successfully connected to database")
	return client, nil
}

// getURI retrieves URI connection with Mongo image
func (db Mongo) getURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?timeoutMS=5000", db.config.Username, db.config.Password, db.config.Host, db.config.Port)
}

// InsertPriceSettings inserts a new document into the PriceSettings collection
func (db Mongo) InsertPriceSettings(settings models.PriceSettings) error {
	// Ensure unique index (only needs to be done once)
	indexModel := mongo.IndexModel{
		Keys: bson.M{"user_id": 1}, // Unique on "user_id" field
		Options: options.Index().
			SetUnique(true),
	}

	_, err := db.collectionClient.Indexes().CreateOne(db.ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index: %s", err.Error())
	}

	result, err := db.collectionClient.InsertOne(db.ctx, settings)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("failed to insert: document already exists")
		} else {
			return fmt.Errorf("failed to insert: %s", err.Error())
		}
	}

	db.logger.Info("update price settings successfully", zap.Any("updated_id", result.InsertedID))
	return err
}

// GetPriceSettings retrieves a document by UserID
func (db Mongo) GetPriceSettings(userID string) (*models.PriceSettings, error) {
	filter := bson.M{"user_id": userID}

	var settings models.PriceSettings
	err := db.collectionClient.FindOne(db.ctx, filter).Decode(&settings)
	if err != nil {
		return nil, fmt.Errorf("failed to get price setting: %s", err.Error())
	}
	db.logger.Info("get price settings successfully", zap.Any("user_id", userID))
	return &settings, nil
}

// PatchPriceSettings updates partial data for user's price settings.
func (db Mongo) PatchPriceSettings(settings models.PriceSettings) error {
	filter := bson.M{"user_id": settings.UserID}

	updates := bson.M{
		"$set": bson.M{
			"vat_included": settings.VatIncluded,
			"margin":       settings.Marginal,
		},
	}
	result, err := db.collectionClient.UpdateOne(db.ctx, filter, updates)
	if err != nil {
		return fmt.Errorf("failed to update price settings: %s", err.Error())
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to update price settings: no matched settings were found")
	}
	db.logger.Info("update price settings successfully", zap.Any("updated_amount", result.MatchedCount))
	return nil
}

// DeletePriceSettings deletes user's price settings.
func (db Mongo) DeletePriceSettings(userID string) error {
	filter := bson.M{"user_id": userID}

	result, err := db.collectionClient.DeleteOne(db.ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete price settings: %s", err.Error())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("failed to delete price settings: no matched settings were found")
	}
	db.logger.Info("delete user price settings successfully", zap.Int64("deleted_amount", result.DeletedCount))
	return nil
}
