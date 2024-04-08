package electric

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

const (
	BASE_URL     string = "https://oomi.fi/wp-json"
	SPOT_PRICE   string = "spot-price"
	GET_V1       string = "v1/get"
	CLIENT_ERROR string = "client"
	SERVER_ERROR string = "server"
)

// Receive request body as struct, beautify it and return as URL string.
// Then call this URL in GET request and decode it
func FetchSpotPrice(requestParameters PriceRequest) (responseData *PriceResponse, errorType string, err error) {
	externalUrl, err := formatRequestParameters(requestParameters)
	if err != nil {
		logger.Logger.Error("failed to format url", zap.Error(err))
		return nil, CLIENT_ERROR, err
	}

	// Make HTTP request to the external source
	resp, err := http.Get(externalUrl)
	if err != nil {
		logger.Logger.Error("failed to fetch data from external source (Oomi)", zap.Error(err))
		return nil, SERVER_ERROR, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil { // Parse []byte to the go struct pointer
		logger.Logger.Error("can not unmarshal JSON", zap.Error(err))
		return nil, SERVER_ERROR, err
	}

	return
}
