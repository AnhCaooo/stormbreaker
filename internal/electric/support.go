package electric

import (
	"fmt"
	"math"
)

var supportedGroups = []string{"hour", "day", "week", "month", "year"}

func formatRequestParameters(requestParameters PriceRequest) (endPoint string, err error) {
	url := fmt.Sprintf("%s/%s/%s", BASE_URL, SPOT_PRICE, GET_V1)

	if !isValidString(requestParameters.StartDate) {
		return "", fmt.Errorf("start date should have value")
	}

	if !isValidString(requestParameters.EndDate) {
		return "", fmt.Errorf("end date should have value")
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

func isValidString(value string) bool {
	return value != ""
}

func isValidGroup(value string) bool {
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
