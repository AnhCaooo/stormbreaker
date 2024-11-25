package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AnhCaooo/go-goods/encode"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// GetPriceSettings retrieves the price settings for specified user
func GetPriceSettings(w http.ResponseWriter, r *http.Request) {
	reqBody, err := encode.DecodeRequest[string](r)
	if err != nil {
		logger.Logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := db.GetPriceSettings(context.TODO(), reqBody)
	if err != nil {
		logger.Logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := encode.EncodeResponse(w, http.StatusOK, settings); err != nil {
		logger.Logger.Error(
			fmt.Sprintf("%s failed to encode response body:", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// CreatePriceSettings creates a new price settings for new user
func CreatePriceSettings(w http.ResponseWriter, r *http.Request) {
	reqBody, err := encode.DecodeRequest[models.PriceSettings](r)
	if err != nil {
		logger.Logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.InsertPriceSettings(context.TODO(), reqBody); err != nil {
		logger.Logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := encode.EncodeResponse(w, http.StatusCreated, reqBody); err != nil {
		logger.Logger.Error(
			fmt.Sprintf("%s failed to encode response body:", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PatchPriceSettings updates the price settings for specified user
func PatchPriceSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// DeletePriceSettings deletes the price settings when user was deleted or removed
func DeletePriceSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
