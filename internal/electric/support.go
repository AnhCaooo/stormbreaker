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
