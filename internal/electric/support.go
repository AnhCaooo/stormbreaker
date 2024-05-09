package electric

import (
	"fmt"
	"math"
	"time"
)

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
	const layout string = "2006-01-02" // this is just the layout of YYYY-MM-DD
	dateTime, err = time.Parse(layout, value)
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
