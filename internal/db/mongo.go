// AnhCao 2024
package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Mongo struct {
	config     *models.Database
	logger     *zap.Logger
	ctx        context.Context
	collection *mongo.Collection
}

func NewMongo(ctx context.Context, config *models.Database, logger *zap.Logger) *Mongo {
	return &Mongo{
		config: config,
		logger: logger,
		ctx:    ctx,
	}
}

// Init to connect to mongo database instance and create collection if it does not exist
func (db *Mongo) EstablishConnection() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db.getURI())
	client, err := mongo.Connect(db.ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	err = client.Ping(db.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err.Error())
	}

	db.collection = client.Database(db.config.Name).Collection(db.config.Collection)

	db.logger.Info("Successfully connected to database")
	return client, nil
}

// getURI retrieves URI connection with Mongo image
func (db Mongo) getURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?timeoutMS=5000", db.config.Username, db.config.Password, db.config.Host, db.config.Port)
}

// todo: unit test
// GetPriceSettings retrieves a document by UserID
func (db Mongo) GetPriceSettings(userID string) (settings *models.PriceSettings, statusCode int, err error) {
	filter := bson.M{"user_id": userID}

	if err = db.collection.FindOne(db.ctx, filter).Decode(settings); err != nil {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to get price setting: %s", err.Error())
		return
	}
	db.logger.Info("get price settings successfully", zap.Any("user_id", userID))
	return settings, http.StatusOK, nil
}

// todo: unit test
// InsertPriceSettings inserts a new document into the PriceSettings collection.
func (db Mongo) InsertPriceSettings(settings models.PriceSettings) (statusCode int, err error) {
	if settings.UserID == "" {
		statusCode = http.StatusUnauthorized
		err = fmt.Errorf("cannot insert un-authenticated document")
		return
	}

	// Ensure unique index (only needs to be done once)
	indexModel := mongo.IndexModel{
		Keys: bson.M{"user_id": 1}, // Unique on "user_id" field
		Options: options.Index().
			SetUnique(true),
	}

	_, err = db.collection.Indexes().CreateOne(db.ctx, indexModel)
	if err != nil {
		statusCode = http.StatusInternalServerError
		err = fmt.Errorf("failed to create index: %s", err.Error())
		return
	}

	result, err := db.collection.InsertOne(db.ctx, settings)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			statusCode = http.StatusConflict
			err = fmt.Errorf("failed to insert: document already exists")
			return
		} else {
			statusCode = http.StatusInternalServerError
			err = fmt.Errorf("failed to insert: %s", err.Error())
			return
		}
	}

	db.logger.Info("update price settings successfully", zap.Any("updated_id", result.InsertedID))
	return http.StatusCreated, err
}

// todo: unit test
// PatchPriceSettings updates partial data for user's price settings.
func (db Mongo) PatchPriceSettings(settings models.PriceSettings) (statusCode int, err error) {
	if settings.UserID == "" {
		statusCode = http.StatusUnauthorized
		err = fmt.Errorf("cannot insert un-authenticated document")
		return
	}

	filter := bson.M{"user_id": settings.UserID}
	updates := bson.M{
		"$set": bson.M{
			"vat_included": settings.VatIncluded,
			"margin":       settings.Marginal,
		},
	}
	result, err := db.collection.UpdateOne(db.ctx, filter, updates)
	if err != nil {
		statusCode = http.StatusInternalServerError
		err = fmt.Errorf("failed to update price settings: %s", err.Error())
		return
	}
	if result.MatchedCount == 0 {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to update price settings: no matched settings were found")
		return
	}
	db.logger.Info("update price settings successfully", zap.Any("updated_amount", result.MatchedCount))
	return http.StatusOK, nil
}

// todo: unit test
// DeletePriceSettings deletes user's price settings.
func (db Mongo) DeletePriceSettings(userID string) (statusCode int, err error) {
	filter := bson.M{"user_id": userID}

	result, err := db.collection.DeleteOne(db.ctx, filter)
	if err != nil {
		statusCode = http.StatusInternalServerError
		err = fmt.Errorf("failed to delete price settings: %s", err.Error())
		return
	}
	if result.DeletedCount == 0 {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to delete price settings: no matched settings were found")
		return
	}
	db.logger.Info("delete user price settings successfully", zap.Int64("deleted_amount", result.DeletedCount))
	return http.StatusOK, nil
}
