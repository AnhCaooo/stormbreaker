package electric

import (
	"fmt"
	"math"
	"time"
)

const DATE_FORMAT string = "2006-01-02" // this is just the layout of YYYY-MM-DD

func formatRequestParameters(requestParameters PriceRequest) (endPoint string, err error) {
	url := fmt.Sprintf("%s/%s/%s", BASE_URL, SPOT_PRICE, GET_V1)

	isValidDateRange, err := isValidDateRange(requestParameters.StartDate, requestParameters.EndDate)
	if !isValidDateRange {
		return "", err
	}

	if !isValidGroup(requestParameters.Group) {
		return "", fmt.Errorf("group should have valid value: 'hour', 'day', 'week', 'month', 'year'")
	}

	if !isValidFloat(requestParameters.Marginal) {
		return "", fmt.Errorf("marginal should have value or equal to 0")
	}

	if !isValidInt(requestParameters.VatIncluded) {
		return "", fmt.Errorf("vatIncluded needs to be value '0' or '1' only")
	}

	if !isValidInt(requestParameters.CompareToLastYear) {
		return "", fmt.Errorf("CompareToLastYear needs to be value '0' or '1' only")
	}

	return fmt.Sprintf("%s?starttime=%s&endtime=%s&margin=%f&group=%s&include_vat%d&compare_to_last_year=%d",
		url, requestParameters.StartDate, requestParameters.EndDate, requestParameters.Marginal,
		requestParameters.Group, requestParameters.VatIncluded, requestParameters.CompareToLastYear,
	), nil
}

func isValidFloat(value float64) bool {
	// Check if value is not NaN and not infinite
	if !math.IsNaN(value) && !math.IsInf(value, 0) {
		return true
	}
	return false
}

func validDate(value string) (dateTime time.Time, err error) {
	dateTime, err = time.Parse(DATE_FORMAT, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("date should have value in correct format 'YYYY-MM-DD'")
	}
	return dateTime, nil
}

func isValidDateRange(startDate, endDate string) (isValid bool, err error) {

	startDateTime, err := validDate(startDate)
	if err != nil {
		return isValid, fmt.Errorf("failed to validate start date. Error: %s", err.Error())
	}

	endDateTime, err := validDate(endDate)
	if err != nil {
		return isValid, fmt.Errorf("failed to validate end date. Error: %s", err.Error())
	}

	if startDateTime == endDateTime || startDateTime.Before(endDateTime) {
		return true, nil
	}
	return false, fmt.Errorf("start date cannot after end date")
}

func isValidGroup(value string) bool {
	var supportedGroups = []string{"hour", "day", "week", "month", "year"}
	for _, group := range supportedGroups {
		if value == group {
			return true
		}
	}
	return false
}

func isValidInt(value int32) bool {
	if value < 0 || value > 1 {
		return false
	}
	return true
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

func getTodayPrices(response PriceResponse) (todayPrice *DailyPrice, err error) {
	filteredPrices := make([]Data, 0)
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

	todayPrice = &DailyPrice{
		Available: pricesAvailable,
		Prices: PriceSeries{
			Name: priceUnit,
			Data: filteredPrices,
		},
	}
	return
}

func getTomorrowPrices(response PriceResponse) (tomorrowPrice *DailyPrice, err error) {
	filteredPrices := make([]Data, 0)
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
		filteredPrices = []Data{}
	}

	tomorrowPrice = &DailyPrice{
		Available: pricesAvailable,
		Prices: PriceSeries{
			Name: priceUnit,
			Data: filteredPrices,
		},
	}
	return
}
