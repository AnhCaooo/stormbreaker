// AnhCao 2024
package db

import (
	"context"
	"net/http"
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
			mockMongo := &Mongo{
				config:     nil,
				logger:     logger,
				ctx:        ctx,
				collection: mt.Coll,
			}

			// Add mock response if expected
			if test.expectedResponse != nil {
				mt.AddMockResponses(test.expectedResponse)
			}

			// Call the function
			statusCode, err := mockMongo.InsertPriceSettings(test.priceSettings)
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

// func TestGetPriceSettings(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Client.Disconnect(context.Background())

// 	tests := []struct {
// 		name string
// 	}{
// 		{name: "get price settings from valid userid"},
// 		{name: "fail to get price settings from invalid userid"},
// 	}

// 	for _, test := range tests {
// 		mt.Run(test.name, func(mt *mtest.T) {

// 		})
// 	}
// }

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
