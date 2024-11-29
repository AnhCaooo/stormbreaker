package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnhCaooo/go-goods/encode"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// Fetch the market spot price of electric in Finland in any times
func (h Handler) PostMarketPrice(w http.ResponseWriter, r *http.Request) {
	reqBody, err := encode.DecodeRequest[models.PriceRequest](r)
	if err != nil {
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	electric := electric.NewElectric(h.logger)
	externalData, errorType, err := electric.FetchSpotPrice(reqBody)
	if err != nil {
		if errorType == models.SERVER_ERROR {
			h.logger.Error(constants.Server, zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := encode.EncodeResponse(w, http.StatusOK, externalData); err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to encode data from external source", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("get market price of electric successfully")
}

// Fetch and return the exchange price for today and tomorrow.
// If tomorrow's price is not available yet, return empty struct.
// Then client needs to show readable information to indicate that data is not available yet.
func (h Handler) GetTodayTomorrowPrice(w http.ResponseWriter, r *http.Request) {
	cacheKey := "today-tomorrow-exchange-price"

	cachePrice, isValid := h.cache.Get(cacheKey)
	if isValid {
		if err := encode.EncodeResponse(w, http.StatusOK, cachePrice); err != nil {
			h.logger.Error(
				fmt.Sprintf("%s failed to encode cache data", constants.Server),
				zap.Error(err),
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.logger.Info("[cache] get today and tomorrow's exchange price successfully")
		return
	}

	electric := electric.NewElectric(h.logger)
	todayTomorrowResponse, err := electric.FetchCurrentSpotPrice(w)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to fetch today and/or tomorrow spot price from external source", constants.Server),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache response to improve performance
	// if tomorrow price is available already, then cache until 23:59
	if todayTomorrowResponse.Tomorrow.Available {
		expiredTime, err := helpers.SetTime(23, 59)
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("%s failed to set expired time for caching", constants.Server),
				zap.Error(err),
			)
			return
		}
		h.cache.SetExpiredAtTime(cacheKey, &todayTomorrowResponse, expiredTime)
		return
	}
	// if tomorrow price is not available and sending request time is before 14:00, then cache until 14:00
	expiredTime, err := helpers.SetTime(14, 00)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("%s failed to set expired time for caching", constants.Server),
			zap.Error(err),
		)
		return
	}
	if time.Now().Before(expiredTime) {
		h.cache.SetExpiredAtTime(cacheKey, &todayTomorrowResponse, expiredTime)
	}
}
