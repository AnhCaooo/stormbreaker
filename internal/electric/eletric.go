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

// get request body for '/market-price/today-tomorrow'
func BuildTodayTomorrowAsBodyRequest() (body PriceRequest, err error) {
	today, tomorrow := getTodayAndTomorrowDateAsString()

	body = PriceRequest{
		StartDate:         today,
		EndDate:           tomorrow,
		Marginal:          0.59, // todo: this field should has default value at the beginning. However, it would be nice to give users have their own customizations and then read from db as it is different between users
		Group:             "hour",
		VatIncluded:       1, // todo: this field by default should be 1. However, it would be nice to give users have their own customizations
		CompareToLastYear: 0,
	}
	return body, nil
}

// todo: any better ways to optimize this code. Go routine?
func MapToTodayTomorrowResponse(data *PriceResponse) (response *TodayTomorrowPrice, err error) {
	todayPrices, err := getTodayPrices(*data)
	if err != nil {
		return nil, err
	}

	tomorrowPrices, err := getTomorrowPrices(*data)
	if err != nil {
		return nil, err
	}

	response = &TodayTomorrowPrice{
		Today:    *todayPrices,
		Tomorrow: *tomorrowPrices,
	}
	return
}
