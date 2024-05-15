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
