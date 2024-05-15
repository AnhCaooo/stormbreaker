package helpers

import (
	"fmt"
	"math"
	"time"
)

const DATE_FORMAT string = "2006-01-02" // this is just the layout of YYYY-MM-DD

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
