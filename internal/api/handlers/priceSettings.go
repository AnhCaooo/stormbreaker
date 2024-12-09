package handlers

import (
	"fmt"
	"net/http"

	"github.com/AnhCaooo/go-goods/encode"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// todo: cache the price settings to improve performance
// GetPriceSettings retrieves the price settings for specified user
func (h Handler) GetPriceSettings(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	settings, statusCode, err := h.mongo.GetPriceSettings(userid)
	if err != nil {
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := encode.EncodeResponse(w, statusCode, settings); err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to encode response body:", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// CreatePriceSettings creates a new price settings for new user
func (h Handler) CreatePriceSettings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	reqBody, err := encode.DecodeRequest[models.PriceSettings](r)
	if err != nil {
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Patch userID from accessToken to price settings struct
	reqBody.UserID = userId
	statusCode, err := h.mongo.InsertPriceSettings(reqBody)
	if err != nil {
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), statusCode)
		return
	}

	response := map[string]string{
		"message": "Operation completed successfully",
	}
	if err := encode.EncodeResponse(w, statusCode, response); err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to encode response body:", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PatchPriceSettings updates the price settings for specified user
func (h Handler) PatchPriceSettings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	reqBody, err := encode.DecodeRequest[models.PriceSettings](r)
	if err != nil {
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Patch userID from accessToken to price settings struct
	reqBody.UserID = userId
	statusCode, err := h.mongo.PatchPriceSettings(reqBody)
	if err != nil {
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), statusCode)
		return
	}

	response := map[string]string{
		"message": "Operation completed successfully",
	}

	if err := encode.EncodeResponse(w, statusCode, response); err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to encode response body:", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// todo: maybe only Admin can perform this action? (to be considered)
// DeletePriceSettings deletes the price settings when user was deleted or removed
func (h Handler) DeletePriceSettings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	statusCode, err := h.mongo.DeletePriceSettings(userId)
	if err != nil {
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), statusCode)
		return
	}

	response := map[string]string{
		"message": "Operation completed successfully",
	}

	if err := encode.EncodeResponse(w, statusCode, response); err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to encode response body:", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
