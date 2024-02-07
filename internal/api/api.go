package api

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

// Fetch the market spot price of electric in Finland
func GetMarketPrice(w http.ResponseWriter, r *http.Request) {
	var reqBody electric.PriceRequest
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.Logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	externalData, err := electric.FetchSpotPrice(reqBody)
	if err != nil {
		logger.Logger.Error("failed to fetch data from external source", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(externalData); err != nil {
		logger.Logger.Error("failed to encode response data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Logger.Debug("stock price of electric", zap.Any("market-price", externalData))
	logger.Logger.Info("get market price of electric successfully")
}
