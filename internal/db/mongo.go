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
	Client     *mongo.Client
	collection *mongo.Collection
}

func NewMongo(ctx context.Context, config *models.Database, logger *zap.Logger) *Mongo {
	return &Mongo{
		config: config,
		logger: logger,
		ctx:    ctx,
	}
}

// EstablishConnection tries to connect to mongo server and create collection if it does not exist
func (db *Mongo) EstablishConnection() (err error) {
	clientOptions := options.Client().ApplyURI(db.getURI())
	db.Client, err = mongo.Connect(db.ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	err = db.Client.Ping(db.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping database: %s", err.Error())
	}

	if err = db.initializeCollection(); err != nil {
		return err
	}
	db.logger.Info("Successfully connected to database")
	return nil
}

func (db *Mongo) initializeCollection() error {
	db.collection = db.Client.Database(db.config.Name).Collection(db.config.Collection)

	// Ensure unique index (only needs to be done once)
	indexModel := mongo.IndexModel{
		Keys: bson.M{"user_id": 1}, // Unique on "user_id" field
		Options: options.Index().
			SetUnique(true),
	}

	_, err := db.collection.Indexes().CreateOne(db.ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index while initialize collection: %s", err.Error())
	}

	return nil
}

// getURI retrieves URI connection with Mongo image
func (db Mongo) getURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?timeoutMS=5000", db.config.Username, db.config.Password, db.config.Host, db.config.Port)
}

// GetPriceSettings retrieves a document by UserID
func (db Mongo) GetPriceSettings(userID string) (settings *models.PriceSettings, statusCode int, err error) {
	if userID == "" {
		statusCode = http.StatusUnauthorized
		err = fmt.Errorf("cannot get price settings from unauthenticated user")
		return
	}
	settings = &models.PriceSettings{}
	filter := bson.M{"user_id": userID}
	if err = db.collection.FindOne(db.ctx, filter).Decode(settings); err != nil {
		settings = nil
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to get price settings: %s", err.Error())
		return
	}
	db.logger.Info("get price settings successfully")
	return settings, http.StatusOK, nil
}

// InsertPriceSettings inserts a new document into the PriceSettings collection.
func (db Mongo) InsertPriceSettings(settings models.PriceSettings) (statusCode int, err error) {
	if settings.UserID == "" {
		statusCode = http.StatusUnauthorized
		err = fmt.Errorf("cannot insert un-authenticated document")
		return
	}

	_, err = db.collection.InsertOne(db.ctx, settings)
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

	db.logger.Info("create new price settings successfully")
	return http.StatusCreated, err
}

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

// DeletePriceSettings deletes user's price settings.
func (db Mongo) DeletePriceSettings(userID string) (statusCode int, err error) {
	if userID == "" {
		statusCode = http.StatusUnauthorized
		err = fmt.Errorf("cannot get price settings from unauthenticated user")
		return
	}
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
