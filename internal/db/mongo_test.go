// AnhCao 2024
package db

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/AnhCaooo/go-goods/log"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap/zapcore"
)

func TestInsertPriceSettings(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	logger := log.InitLogger(zapcore.InfoLevel)
	ctx := context.TODO()

	tests := []struct {
		name               string
		priceSettings      models.PriceSettings
		expectedResponse   bson.D
		expectedStatusCode int
		expectedError      string
	}{
		{
			name: "successful operation/create new price settings with valid struct",
			priceSettings: models.PriceSettings{
				UserID:      "12345",
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedResponse:   mtest.CreateSuccessResponse(),
			expectedStatusCode: http.StatusCreated,
			expectedError:      "",
		},
		{
			name: "general error/something went wrong while create new price settings",
			priceSettings: models.PriceSettings{
				UserID:      "12345",
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedResponse: mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    12345, // Some other error code
					Message: "some database error",
				},
			),
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "failed to insert: some database error",
		},
		{
			name: "unauthorized insertion/create new price settings without userid",
			priceSettings: models.PriceSettings{
				UserID:      "",
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedResponse:   nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      "cannot insert un-authenticated document",
		},
		{
			name: "duplicated error/create duplicated price settings (same userid)",
			priceSettings: models.PriceSettings{
				UserID:      "12345",
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedResponse: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index: 0,
				Code:  11000, // Duplicate key error code
			}),
			expectedStatusCode: http.StatusConflict,
			expectedError:      "failed to insert: document already exists",
		},
	}

	for _, test := range tests {
		mt.Run(test.name, func(mt *mtest.T) {
			db := NewMongo(ctx, nil, logger)
			db.collection = mt.Coll

			// Add mock response if expected
			if test.expectedResponse != nil {
				mt.AddMockResponses(test.expectedResponse)
			}

			// Call the function
			statusCode, err := db.InsertPriceSettings(test.priceSettings)
			// Validate error
			if err != nil && err.Error() != test.expectedError {
				t.Errorf("got %q, wanted %q", err.Error(), test.expectedError)
			}

			// Validate status code
			if err == nil && statusCode != test.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", test.expectedStatusCode, statusCode)
			}
		})
	}

}

func TestGetPriceSettings(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	logger := log.InitLogger(zapcore.DebugLevel)
	ctx := context.TODO()

	tests := []struct {
		name               string
		userID             string
		mockResponse       bson.D
		expectedSettings   models.PriceSettings
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:   "successful operation/valid user ID and price settings found",
			userID: "12345",
			mockResponse: mtest.CreateCursorResponse(1, "test.price_settings", mtest.FirstBatch, bson.D{
				{Key: "user_id", Value: "12345"},
				{Key: "margin", Value: 0.59},
				{Key: "vat_included", Value: true},
			}),
			expectedSettings: models.PriceSettings{
				UserID:      "12345",
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedStatusCode: http.StatusOK,
			expectedError:      "",
		},
		{
			name:               "unauthenticated user/empty user ID",
			userID:             "",
			mockResponse:       nil,
			expectedSettings:   models.PriceSettings{},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      "cannot get price settings from unauthenticated user",
		},
		{
			name:               "price settings not found",
			userID:             "54321",
			mockResponse:       mtest.CreateCursorResponse(0, "test.price_settings", mtest.FirstBatch),
			expectedSettings:   models.PriceSettings{},
			expectedStatusCode: http.StatusNotFound,
			expectedError:      "failed to get price settings: mongo: no documents in result",
		},
	}

	for _, test := range tests {
		mt.Run(test.name, func(mt *mtest.T) {
			db := NewMongo(ctx, nil, logger)
			db.collection = mt.Coll

			// Set up mock responses if applicable
			if test.mockResponse != nil {
				mt.AddMockResponses(test.mockResponse)
			}

			// Call the function being tested
			settings, statusCode, err := db.GetPriceSettings(test.userID)
			// Validate error
			if err != nil && err.Error() != test.expectedError {
				t.Errorf("unexpected error: got %q, want %q", err.Error(), test.expectedError)
			} else if err == nil && test.expectedError != "" {
				t.Errorf("expected error %q but got nil", test.expectedError)
			}

			// Validate status code
			if statusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", statusCode, test.expectedStatusCode)
			}

			defaultSettings := models.PriceSettings{}
			// Validate empty settings (if any)
			if test.expectedSettings == defaultSettings && settings != nil {
				t.Errorf("expected empty settings: %#v but got: %#v", test.expectedSettings, &settings)
			}
			// Validate settings (if any)
			if test.expectedSettings != defaultSettings && settings != nil {
				if !reflect.DeepEqual(settings, &test.expectedSettings) {
					t.Errorf("expected settings: %#v but got: %#v", test.expectedSettings, settings)
				}
			}
		})
	}

}

func TestPatchPriceSettings(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	logger := log.InitLogger(zapcore.DebugLevel)
	ctx := context.TODO()

	tests := []struct {
		name               string
		priceSettings      models.PriceSettings
		mockResponse       bson.D
		expectedStatusCode int
		expectedError      string
	}{
		{
			name: "successful update: valid user ID and updated fields",
			priceSettings: models.PriceSettings{
				UserID:      "12345",
				Marginal:    0.75,
				VatIncluded: true,
			},
			mockResponse: bson.D{
				{Key: "n", Value: 1},         // Number of matched documents
				{Key: "nModified", Value: 1}, // Number of modified documents
				{Key: "ok", Value: 1},
			},
			expectedStatusCode: http.StatusOK,
			expectedError:      "",
		},
		{
			name: "unauthorized operation: empty user ID",
			priceSettings: models.PriceSettings{
				UserID:      "",
				Marginal:    0.75,
				VatIncluded: true,
			},
			mockResponse:       nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      "cannot insert un-authenticated document",
		},
		{
			name: "not found: no matched settings for the given user ID",
			priceSettings: models.PriceSettings{
				UserID:      "67890",
				Marginal:    0.85,
				VatIncluded: false,
			},
			mockResponse: bson.D{
				{Key: "ok", Value: 1},
				{Key: "nModified", Value: 0},
				{Key: "matchedCount", Value: 0},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError:      "failed to update price settings: no matched settings were found",
		},
		{
			name: "internal server error: database failure",
			priceSettings: models.PriceSettings{
				UserID:      "12345",
				Marginal:    0.75,
				VatIncluded: true,
			},
			mockResponse: mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    12345,
				Message: "some database error",
			}),
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "failed to update price settings: some database error",
		},
	}

	for _, test := range tests {
		mt.Run(test.name, func(mt *mtest.T) {
			// Setup MongoDB mock instance
			db := NewMongo(ctx, nil, logger)
			db.collection = mt.Coll

			// Set up mock responses
			mt.AddMockResponses(test.mockResponse)

			// Call the PatchPriceSettings function
			statusCode, err := db.PatchPriceSettings(test.priceSettings)

			// Validate error
			if test.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", test.expectedError)
				} else if err.Error() != test.expectedError {
					t.Errorf("unexpected error: got %q, want %q", err.Error(), test.expectedError)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Validate status code
			if statusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", statusCode, test.expectedStatusCode)
			}
		})

	}
}

func TestDeletePriceSettings(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	logger := log.InitLogger(zapcore.DebugLevel)
	ctx := context.TODO()

	tests := []struct {
		name               string
		userID             string
		mockResponse       bson.D
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:   "successful deletion",
			userID: "12345",
			mockResponse: bson.D{
				{Key: "ok", Value: 1},
				{Key: "n", Value: 1}, // Number of documents deleted
			},
			expectedStatusCode: http.StatusOK,
			expectedError:      "",
		},
		{
			name:               "empty user ID",
			userID:             "",
			mockResponse:       nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      "cannot get price settings from unauthenticated user",
		},
		{
			name:   "no matched documents: user ID not found",
			userID: "99999",
			mockResponse: bson.D{
				{Key: "ok", Value: 1},
				{Key: "n", Value: 0}, // No documents deleted
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError:      "failed to delete price settings: no matched settings were found",
		},
	}

	for _, test := range tests {
		mt.Run(test.name, func(mt *mtest.T) {
			db := NewMongo(ctx, nil, logger)
			db.collection = mt.Coll

			if test.mockResponse != nil {
				mt.AddMockResponses(test.mockResponse)
			}
			// Call the function being tested
			statusCode, err := db.DeletePriceSettings(test.userID)

			// Validate error
			if test.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", test.expectedError)
				} else if err.Error() != test.expectedError {
					t.Errorf("unexpected error: got %q, want %q", err.Error(), test.expectedError)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Validate status code
			if statusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", statusCode, test.expectedStatusCode)
			}
		})
	}
}
