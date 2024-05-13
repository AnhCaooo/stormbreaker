package api

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

// Fetch the market spot price of electric in Finland in any times
func PostMarketPrice(w http.ResponseWriter, r *http.Request) {
	var reqBody electric.PriceRequest
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.Logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	externalData, errorType, err := electric.FetchSpotPrice(reqBody)
	if err != nil {
		if errorType == electric.SERVER_ERROR {
			logger.Logger.Error("failed to fetch data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Logger.Error("failed to fetch data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")

}
