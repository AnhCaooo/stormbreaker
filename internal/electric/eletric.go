// AnhCao 2024
package electric

import (
	"fmt"
	"net/http"

	"github.com/AnhCaooo/go-goods/encode"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

type Electric struct {
	logger        *zap.Logger
	mongo         *db.Mongo
	userId        string
	priceSettings *models.PriceSettings
}

func NewElectric(logger *zap.Logger, mongo *db.Mongo, userId string, priceSettings *models.PriceSettings) *Electric {
	if mongo == nil {
		logger.Warn("MongoDB client is nil, using mock or no-op database")
	}

	return &Electric{
		logger:        logger,
		mongo:         mongo,
		userId:        userId,
		priceSettings: priceSettings,
	}
}

// Receive request body as struct, beautify it and return as URL string.
// Then call this URL in GET request and decode it
func (e Electric) FetchSpotPrice(requestParameters *models.PriceRequest) (responseData *models.PriceResponse, statusCode int, err error) {
	var settings *models.PriceSettings = e.priceSettings
	if e.mongo == nil || e.priceSettings == nil || e.userId == "stormbreaker" {
		e.logger.Debug("load default price settings")
		settings = e.getDefaultPriceSettings()
	}

	externalUrl, err := helpers.FormatMarketPricePostReqParameters(requestParameters, settings)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Make HTTP request to the external source
	resp, err := http.Get(externalUrl)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to fetch data from external source (Oomi): %s", err.Error())
	}
	defer resp.Body.Close()

	responseData, err = encode.DecodeResponse[*models.PriceResponse](resp)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	statusCode = http.StatusOK
	return
}

// FetchCurrentSpotPrice retrieves the current spot price for today and tomorrow,
// maps the data to a response structure, and writes the response to the provided
// http.ResponseWriter. It returns the mapped response and any error encountered.
// Depending on the time sending request, there could be tomorrow's price come along with today's price.
// In practice, tomorrow's price would be available around 3pm (Finnish time) everyday.
func (e Electric) FetchCurrentSpotPrice(w http.ResponseWriter) (todayTomorrowResponse *models.TodayTomorrowPrice, err error) {
	reqBody := e.BuildTodayTomorrowRequestPayload()
	todayTomorrowPrice, _, err := e.FetchSpotPrice(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s failed to fetch data: %s", constants.Server, err.Error())
	}

	todayTomorrowResponse, err = helpers.MapToTodayTomorrowResponse(todayTomorrowPrice)
	if err != nil {
		return nil, fmt.Errorf("%s failed to map to informative struct data: %s", constants.Server, err.Error())
	}

	if err := encode.EncodeResponse(w, http.StatusOK, todayTomorrowResponse); err != nil {
		return nil, fmt.Errorf("%s failed to encode response data: %s", constants.Server, err.Error())
	}

	e.logger.Info("[from external source] get today and tomorrow's exchange price successfully")
	return
}

// return as request body with date of today and data of tomorrow.
// Usage: get request body for '/market-price/today-tomorrow'
func (e Electric) BuildTodayTomorrowRequestPayload() *models.PriceRequest {
	today, tomorrow := helpers.GetTodayAndTomorrowDateAsString()

	return &models.PriceRequest{
		StartDate:         today,
		EndDate:           tomorrow,
		Group:             "hour",
		CompareToLastYear: 0,
	}
}

// GetDefaultPriceSettings returns a default values in case the service cannot get the price settings from database.s
func (e Electric) getDefaultPriceSettings() *models.PriceSettings {
	return &models.PriceSettings{
		UserID:      e.userId,
		Marginal:    0.0,
		VatIncluded: false,
	}
}
