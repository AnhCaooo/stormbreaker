package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnhCaooo/go-goods/encode"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// todo: cache the price settings to improve performance
// GetPriceSettings retrieves the price settings for specific user
//
//	@Summary		Retrieves the price settings for specific user
//	@Description	retrieves the price settings for specific user by identify through 'access token'.
//	@Tags			price-settings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.PriceSettings
//	@Failure		400	{object}	string "Invalid request"
//	@Failure		401	{object}	string "Unauthenticated/Unauthorized"
//	@Failure		500	{object}	string "Various reasons: cannot fetch price from 3rd party, failed to read settings from db, etc."
//	@Router			/v1/price-settings [get]
func (h Handler) GetPriceSettings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	cacheKey := fmt.Sprintf("%s-price-settings", userId)

	cacheSettings, isValid := h.cache.Get(cacheKey)
	if isValid {
		if err := encode.EncodeResponse(w, http.StatusOK, cacheSettings); err != nil {
			h.logger.Error(
				fmt.Sprintf("%s failed to encode cache data", constants.Server),
				zap.Error(err),
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.logger.Info("[cache] get price settings successfully")
		return
	}

	settings, statusCode, err := h.mongo.GetPriceSettings(userId)
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

	// cache price settings and keep for 1 hours
	h.cache.SetExpiredAfterTimePeriod(cacheKey, &settings, time.Hour*time.Duration(1))

}

// CreatePriceSettings creates a new price settings for user
//
//	@Summary		Creates a new price settings for user
//	@Description	Creates a new price settings for new user by identify through 'access token'.
//	@Tags			price-settings
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.PriceSettings	true	"user price settings"
//	@Success		200	{object}	string
//	@Failure		400	{object}	string "Invalid request"
//	@Failure		401	{object}	string "Unauthenticated/Unauthorized"
//	@Failure		403	{object}	string "Forbidden"
//	@Failure		404	{object}	string "Settings not found"
//	@Failure		409	{object}	string "Settings exist already"
//	@Failure		500	{object}	string "Various reasons: cannot fetch price from 3rd party, failed to read settings from db, etc."
//	@Router			/v1/price-settings [post]
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

	if reqBody.UserID != "" && reqBody.UserID != userId {
		err = fmt.Errorf("given `user_id` %s is different from `user_id` in `access_token`", reqBody.UserID)
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
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
//
//	@Summary		Updates the price settings for specific user
//	@Description	Updates the price settings for specific user by identify through 'access token'.
//	@Tags			price-settings
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.PriceSettings	true	"user price settings"
//	@Success		200	{object}	string
//	@Failure		400	{object}	string "Invalid request"
//	@Failure		401	{object}	string "Unauthenticated/Unauthorized"
//	@Failure		403	{object}	string "Forbidden"
//	@Failure		404	{object}	string "Settings not found"
//	@Failure		500	{object}	string "Various reasons: cannot fetch price from 3rd party, failed to read settings from db, etc."
//	@Router			/v1/price-settings [patch]
func (h Handler) PatchPriceSettings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	cacheKey := fmt.Sprintf("%s-price-settings", userId)

	reqBody, err := encode.DecodeRequest[models.PriceSettings](r)
	if err != nil {
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserID != "" && reqBody.UserID != userId {
		err = fmt.Errorf("given `user_id` %s is different from `user_id` in `access_token`", reqBody.UserID)
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
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

	h.cache.Delete(cacheKey)
}

// todo: maybe only Admin can perform this action? (to be considered)
// DeletePriceSettings deletes the price settings when user was deleted or removed
//
//	@Summary		Deletes the price settings for specific user
//	@Description	Deletes the price settings for specific user by identify through 'access token'.
//	@Tags			price-settings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Failure		400	{object}	string "Invalid request"
//	@Failure		401	{object}	string "Unauthenticated/Unauthorized"
//	@Failure		404	{object}	string "Settings not found"
//	@Failure		500	{object}	string "Various reasons: cannot fetch price from 3rd party, failed to read settings from db, etc."
//	@Router			/v1/price-settings [delete]
func (h Handler) DeletePriceSettings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	cacheKey := fmt.Sprintf("%s-price-settings", userId)

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
	h.cache.Delete(cacheKey)

}
