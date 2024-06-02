package electric

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

// Receive request body as struct, beautify it and return as URL string.
// Then call this URL in GET request and decode it
func FetchSpotPrice(requestParameters models.PriceRequest) (responseData *models.PriceResponse, errorType string, err error) {
	externalUrl, err := helpers.FormatRequestParameters(requestParameters)
	if err != nil {
		logger.Logger.Error("failed to format url", zap.Error(err))
		return nil, models.CLIENT_ERROR, err
	}

	// Make HTTP request to the external source
	resp, err := http.Get(externalUrl)
	if err != nil {
		logger.Logger.Error("failed to fetch data from external source (Oomi)", zap.Error(err))
		return nil, models.SERVER_ERROR, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil { // Parse []byte to the go struct pointer
		logger.Logger.Error("can not unmarshal JSON", zap.Error(err))
		return nil, models.SERVER_ERROR, err
	}

	return
}

// fetch and return current spot price. Depending on the time sending request,
// there could be tomorrow's price come along with today's price.
// In practice, tomorrow's price would be available around 2pm-4pm everyday
func FetchCurrentSpotPrice(w http.ResponseWriter) *models.TodayTomorrowPrice {
	var reqBody, err = helpers.BuildTodayTomorrowAsBodyRequest()
	if err != nil {
		logger.Logger.Error("[server error] failed to build request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	todayTomorrowPrice, errorType, err := FetchSpotPrice(reqBody)
	if err != nil {
		if errorType == models.SERVER_ERROR {
			logger.Logger.Error("[server error] failed to fetch data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		logger.Logger.Error("[request error] failed to fetch data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	todayTomorrowResponse, err := helpers.MapToTodayTomorrowResponse(todayTomorrowPrice)
	if err != nil {
		logger.Logger.Error("[server error] failed to map to informative struct data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	if err := json.NewEncoder(w).Encode(todayTomorrowResponse); err != nil {
		logger.Logger.Error("failed to encode response data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	logger.Logger.Info("get today and tomorrow's exchange price successfully")

	return todayTomorrowResponse
}
