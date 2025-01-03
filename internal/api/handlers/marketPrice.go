package handlers

import (
	"fmt"
	"net/http"

	"github.com/AnhCaooo/go-goods/encode"
	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// PostMarketPrice fetches the market spot price of electric in Finland in any times
//
//	@Summary		Retrieves the market price
//	@Description	Fetch the market spot price of electric in Finland in any times
//	@Tags			market-price
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.PriceRequest	true	"Criteria for getting market spot price"
//	@Success		200	{object}	models.PriceResponse
//	@Failure		400	{string}	string "Invalid request"
//	@Failure		401	{string}	string "Unauthenticated/Unauthorized"
//	@Failure		500	{string}	string "Various reasons: cannot fetch price from 3rd party, failed to read settings from db, etc."
//	@Router			/v1/market-price [post]
func (h Handler) PostMarketPrice(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	reqBody, err := encode.DecodeRequest[models.PriceRequest](r)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s", h.workerID, constants.Client), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, _, err := h.LoadPriceSettings(userID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s", h.workerID, constants.Server), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	electric := electric.NewElectric(h.logger, h.mongo, userID, settings)
	externalData, statusCode, err := electric.FetchSpotPrice(&reqBody)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s", h.workerID, constants.Server), zap.Error(err))
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := encode.EncodeResponse(w, statusCode, externalData); err != nil {
		h.logger.Error(
			fmt.Sprintf("[worker_%d] %s failed to encode data from external source", h.workerID, constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info(fmt.Sprintf("[worker_%d] got market price of electric successfully", h.workerID))
}

// GetTodayTomorrowPrice returns the exchange price for today and tomorrow.
// If tomorrow's price is not available yet, return empty struct.
// Then client (Web, mobile) needs to show readable information to indicate that data is not available yet.
//
//	@Summary		Retrieves the market price for today and tomorrow
//	@Description	Returns the exchange price for today and tomorrow.
//	@Description	If tomorrow price is not available yet, return empty struct.
//	@Description	Then client needs to show readable information to indicate that data is not available yet.
//	@Tags			market-price
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.TodayTomorrowPrice
//	@Failure		401	{string}	string "Unauthenticated/Unauthorized"
//	@Failure		500	{string}	string "Various reasons: cannot fetch price from 3rd party, failed to read settings from db, etc."
//	@Router			/v1/market-price/today-tomorrow [get]
func (h Handler) GetTodayTomorrowPrice(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	settings, _, err := h.LoadPriceSettings(userID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s", h.workerID, constants.Server), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cachePrice, isValid := h.cache.Get(cache.PlainTodayTomorrowPricesKey)
	if isValid {
		pricesMessage, err := helpers.MapInterfaceToStruct[models.NewPricesMessage](cachePrice)
		if err != nil {
			h.logger.Error(fmt.Sprintf("[worker_%d] [cache] failed to cast cache data to NewPricesMessage", h.workerID))
			http.Error(w, "Failed to cast cache data to NewPricesMessage", http.StatusInternalServerError)
			return
		}

		// map the price settings with plain current spot price
		todayTomorrowPrices := helpers.MapPriceSettingsWithTodayTomorrowSpotPrice(settings, &pricesMessage.Data)
		if err := encode.EncodeResponse(w, http.StatusOK, todayTomorrowPrices); err != nil {
			h.logger.Error(
				fmt.Sprintf("[worker_%d] %s failed to encode cache data", h.workerID, constants.Server),
				zap.Error(err),
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.logger.Info(fmt.Sprintf("[worker_%d] [cache] get today and tomorrow's exchange price successfully", h.workerID))
		return
	}

	electric := electric.NewElectric(h.logger, h.mongo, userID, settings)
	_, err = electric.FetchCurrentSpotPrice(w)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("[worker_%d] %s failed to fetch today and/or tomorrow spot price from external source", h.workerID, constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// // Cache response to improve performance
	// // if tomorrow price is available already, then cache until 23:59
	// if todayTomorrowResponse.Tomorrow.Available {
	// 	expiredTime, err := helpers.SetTime(23, 59)
	// 	if err != nil {
	// 		h.logger.Error(
	// 			fmt.Sprintf("[worker_%d] %s failed to set expired time for caching", h.workerID, constants.Server),
	// 			zap.Error(err),
	// 		)
	// 		return
	// 	}
	// 	h.cache.SetExpiredAtTime(cacheKey, &todayTomorrowResponse, expiredTime)
	// 	return
	// }
	// // if tomorrow price is not available and sending request time is before 14:00, then cache until 14:00
	// expiredTime, err := helpers.SetTime(14, 00)
	// if err != nil {
	// 	h.logger.Error(
	// 		fmt.Sprintf("[worker_%d] %s failed to set expired time for caching", h.workerID, constants.Server),
	// 		zap.Error(err),
	// 	)
	// 	return
	// }
	// if time.Now().Before(expiredTime) {
	// 	h.cache.SetExpiredAtTime(cacheKey, &todayTomorrowResponse, expiredTime)
	// }
}
