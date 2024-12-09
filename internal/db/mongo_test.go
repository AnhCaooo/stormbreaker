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
				t.Errorf("got expected status code %d, got %d", test.expectedStatusCode, statusCode)
			}
		})
	}

}

func TestGetPriceSettings(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	logger := log.InitLogger(zapcore.DebugLevel)
	ctx := context.Background()

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

// func TestPatchPriceSettings(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Client.Disconnect(context.Background())

// 	tests := []struct {
// 		name string
// 	}{
// 		{name: "update price settings with valid struct"},
// 		{name: "update price settings with invalid userid (userid is empty)"},
// 		{name: "update price settings with invalid userid (userid does not exist)"},
// 		{name: "update price settings with invalid struct"},
// 		{name: "update price settings with same struct"},
// 	}

// 	for _, test := range tests {
// 		mt.Run(test.name, func(mt *mtest.T) {

// 		})
// 	}

// }

// func TestDeletePriceSettings(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Client.Disconnect(context.Background())

// 	tests := []struct {
// 		name string
// 	}{
// 		{name: "delete price settings with valid userid"},
// 		{name: "delete non-existing price settings"},
// 		{name: "delete price settings with invalid userid"},
// 	}

// 	for _, test := range tests {
// 		mt.Run(test.name, func(mt *mtest.T) {

// 		})
// 	}

// }
