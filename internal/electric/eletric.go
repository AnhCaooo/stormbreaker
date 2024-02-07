package electric

import (
	"encoding/json"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

const (
	BASE_URL   = "https://oomi.fi/wp-json"
	SPOT_PRICE = "spot-price"
	GET_V1     = "v1/get"
)

// Receive request body as struct, beautify it and return as URL string.
// Then call this URL in GET request and decode it
func FetchSpotPrice(requestParameters PriceRequest) (responseData *PriceResponse, err error) {
	externalUrl, err := formatRequestParameters(requestParameters)
	if err != nil {
		logger.Logger.Error("failed to format Url", zap.Error(err))
		return nil, err
	}

	// Make HTTP request to the external source
	resp, err := http.Get(externalUrl)
	if err != nil {
		logger.Logger.Error("failed to fetch data from external source (Oomi)", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil { // Parse []byte to the go struct pointer
		logger.Logger.Error("can not unmarshal JSON", zap.Error(err))
		return nil, err
	}

	return
}
