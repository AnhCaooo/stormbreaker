// AnhCao 2024
package helpers

import (
	"fmt"
	"time"

	"github.com/AnhCaooo/stormbreaker/internal/models"
)

// receives 'requestParameters' struct and price settings. Then return appropriate endpoint url
func FormatMarketPricePostReqParameters(requestParameters *models.PriceRequest, settings *models.PriceSettings) (endPoint string, err error) {
	url := fmt.Sprintf("%s/%s/%s", models.BASE_URL, models.SPOT_PRICE, models.GET_V1)

	isValidDateRange, err := isValidDateRange(requestParameters.StartDate, requestParameters.EndDate)
	if !isValidDateRange {
		return "", err
	}

	if !isValidGroup(requestParameters.Group) {
		return "", fmt.Errorf("group should have valid value: 'hour', 'day', 'week', 'month', 'year'")
	}

	if !isValidFloat(settings.Marginal) {
		return "", fmt.Errorf("marginal should have float value or equal to 0")
	}

	if !isValidInt(parseVatIncludedFromBoolToInt32(settings.VatIncluded)) {
		return "", fmt.Errorf("vatIncluded needs to be value '0' or '1' only")
	}

	if !isValidInt(requestParameters.CompareToLastYear) {
		return "", fmt.Errorf("compareToLastYear needs to be value '0' or '1' only")
	}

	return fmt.Sprintf("%s?starttime=%s&endtime=%s&margin=%f&group=%s&include_vat=%d&compare_to_last_year=%d",
		url, requestParameters.StartDate, requestParameters.EndDate, settings.Marginal,
		requestParameters.Group, parseVatIncludedFromBoolToInt32(settings.VatIncluded), requestParameters.CompareToLastYear,
	), nil
}

// receives price's response and map it to `TodayTomorrowPrice` 's struct
func MapToTodayTomorrowResponse(data *models.PriceResponse) (response *models.TodayTomorrowPrice, err error) {
	todayPrices, err := getTodayPrices(*data)
	if err != nil {
		return nil, err
	}

	tomorrowPrices, err := getTomorrowPrices(*data)
	if err != nil {
		return nil, err
	}

	response = &models.TodayTomorrowPrice{
		Today:    *todayPrices,
		Tomorrow: *tomorrowPrices,
	}
	return
}

// GetTodayAndTomorrowDateAsString returns date of today and tomorrow
func GetTodayAndTomorrowDateAsString() (todayDate, tomorrowDate string) {
	// Get today's date
	today := time.Now()
	// Get tomorrow's date by adding one day
	tomorrow := today.AddDate(0, 0, 1)

	// convert Date to string
	todayDate = today.Format(DATE_FORMAT)
	tomorrowDate = tomorrow.Format(DATE_FORMAT)

	return
}

func parseVatIncludedFromBoolToInt32(included bool) int32 {
	if included {
		return 1
	}
	return 0
}

func getTodayPrices(response models.PriceResponse) (todayPrice *models.DailyPrice, err error) {
	filteredPrices := make([]models.Data, 0)
	pricesAvailable := false

	priceUnit := response.Data.Series[0].Name
	pricesData := response.Data.Series[0].Data
	if len(pricesData) == 0 {
		return nil, fmt.Errorf("failed to get price for today from the response: %v", response)
	}

	for _, price := range pricesData {

		if price.IsToday {
			filteredPrices = append(filteredPrices, price)
		}
	}

	if len(filteredPrices) != 24 {
		return nil, fmt.Errorf("the amount of price per hour exceed 24. Its length is %d", len(filteredPrices))
	} else {
		pricesAvailable = true
	}

	todayPrice = &models.DailyPrice{
		Available: pricesAvailable,
		Prices: models.PriceSeries{
			Name: priceUnit,
			Data: filteredPrices,
		},
	}
	return
}

func getTomorrowPrices(response models.PriceResponse) (tomorrowPrice *models.DailyPrice, err error) {
	filteredPrices := make([]models.Data, 0)
	pricesAvailable := false

	priceUnit := response.Data.Series[0].Name
	pricesData := response.Data.Series[0].Data

	if len(pricesData) == 0 {
		return nil, fmt.Errorf("failed to get price for tomorrow from the response: %v", response)
	}

	for _, price := range pricesData {
		if !price.IsToday {
			filteredPrices = append(filteredPrices, price)
		}
	}

	if len(filteredPrices) == 24 {
		pricesAvailable = true
	} else {
		// clear the filtered prices so that client will not get confused why there is only 1 price at 00:00. This is legacy from external source
		// ? At least of now, ignore and show empty data. AnhC - 15th May 2024
		filteredPrices = []models.Data{}
	}

	tomorrowPrice = &models.DailyPrice{
		Available: pricesAvailable,
		Prices: models.PriceSeries{
			Name: priceUnit,
			Data: filteredPrices,
		},
	}
	return
}
