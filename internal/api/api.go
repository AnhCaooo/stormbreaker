// AnhCao 2024
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// Ping the connection to the server
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// Fetch the market spot price of electric in Finland in any times
func PostMarketPrice(w http.ResponseWriter, r *http.Request) {
	reqBody, err := helpers.DecodeRequest[models.PriceRequest](r)
	if err != nil {
		logger.Logger.Error("[request error] failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	externalData, errorType, err := electric.FetchSpotPrice(reqBody)
	if err != nil {
		if errorType == models.SERVER_ERROR {
			logger.Logger.Error("[server] failed to fetch data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Logger.Error("[request error] failed to fetch data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := helpers.EncodeResponse(w, http.StatusOK, externalData); err != nil {
		logger.Logger.Error("[server] failed to encode external response data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("get market price of electric successfully")
}

// Fetch and return the exchange price for today and tomorrow.
// If tomorrow's price is not available yet, return empty struct.
// Then client needs to show readable information to indicate that data is not available yet.
func GetTodayTomorrowPrice(w http.ResponseWriter, r *http.Request) {
	cacheKey := "today-tomorrow-exchange-price"

	cachePrice, isValid := cache.Cache.Get(cacheKey)
	if isValid {
		if err := helpers.EncodeResponse(w, http.StatusOK, cachePrice); err != nil {
			logger.Logger.Error("[server] failed to encode cache data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Logger.Info("[cache] get today and tomorrow's exchange price successfully")
		return
	}

	todayTomorrowResponse, err := electric.FetchCurrentSpotPrice(w)
	if err != nil {
		logger.Logger.Error("[server] failed to fetch today and/or tomorrow spot price from external source", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache response to improve performance
	// if tomorrow price is available already, then cache until 23:59
	if todayTomorrowResponse.Tomorrow.Available {
		expiredTime, err := helpers.SetTime(23, 59)
		if err != nil {
			logger.Logger.Error("[server] failed to set expired time for caching.", zap.Error(err))
			return
		}
		cache.Cache.SetExpiredAtTime(cacheKey, &todayTomorrowResponse, expiredTime)
		return
	}
	// if tomorrow price is not available and sending request time is before 14:00, then cache until 14:00
	expiredTime, err := helpers.SetTime(14, 00)
	if err != nil {
		logger.Logger.Error("[server] failed to set expired time for caching.", zap.Error(err))
		return
	}
	if time.Now().Before(expiredTime) {
		cache.Cache.SetExpiredAtTime(cacheKey, &todayTomorrowResponse, expiredTime)
	}
}
