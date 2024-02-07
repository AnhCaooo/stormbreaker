package api

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	externalData, err := electric.FetchSpotPrice()
	if err != nil {
		logger.Logger.Error("failed to fetch data from external source", zap.Error(err))
		http.Error(w, "failed to fetch data from external source", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(externalData); err != nil {
		logger.Logger.Error("failed to encode response data", zap.Error(err))
		http.Error(w, "failed to encode response data", http.StatusInternalServerError)
		return
	}
}
