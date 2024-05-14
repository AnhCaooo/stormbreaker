package electric

import (
	"encoding/json"
	"net/http"
	"time"

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

func getTodayAndTomorrowDateAsString() (todayDate string, tomorrowDate string) {
	// Get current time
	now := time.Now()

	// Get year, month, and day components
	year, month, day := now.Date()

	// Get today's date
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	// Get tomorrow's date by adding one day
	tomorrow := today.AddDate(0, 0, 1)

	todayDate = today.Format(DATE_FORMAT)
	tomorrowDate = tomorrow.Format(DATE_FORMAT)
	return
}
