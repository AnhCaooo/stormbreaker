package api

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// Fetch the market spot price of electric in Finland in any times
func PostMarketPrice(w http.ResponseWriter, r *http.Request) {
	var reqBody models.PriceRequest

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.Logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	externalData, errorType, err := electric.FetchSpotPrice(reqBody)
	if err != nil {
		if errorType == models.SERVER_ERROR {
			logger.Logger.Error("[server error] failed to fetch data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Logger.Error("[request error] failed to fetch data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(externalData); err != nil {
		logger.Logger.Error("failed to encode response data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Logger.Info("get market price of electric successfully")
}

// Fetch and return the exchange price for today and tomorrow.
// If tomorrow's price is not available yet, return empty struct.
// Then client needs to show readable information to indicate that data is not available yet.
func GetTodayTomorrowPrice(w http.ResponseWriter, r *http.Request) {
	cacheKey := "today-tomorrow-exchange-price"
	expiredAtHour := 14
	w.Header().Set("Content-Type", "application/json")

	cachePrice, isValid := cache.Cache.Get(cacheKey)
	if isValid {
		if err := json.NewEncoder(w).Encode(cachePrice); err != nil {
			logger.Logger.Error("[server error] failed to encode cache data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Logger.Info("get today and tomorrow's exchange price successfully")
		return
	}

	todayTomorrowResponse := electric.FetchCurrentSpotPrice(w)
	cache.Cache.SetExpiredAt(cacheKey, &todayTomorrowResponse, expiredAtHour)
}
