// AnhCao 2024
package helpers

import (
	"fmt"
	"time"
)

// return current date in Helsinki time with specific hour and minute.
// Use case examples: set cache to expired at specific time
func SetTime(hour int, minute int) (time.Time, error) {
	// get timezone
	location, err := loadHelsinkiLocation()
	if err != nil {
		return time.Now(), err
	}

	// Get current time in Finnish time
	now := time.Now().In(location)

	// Get year, month, and day components
	year, month, day := now.Date()
	settingTime := time.Date(year, month, day, hour, minute, 0, 0, location)
	// return as UTC
	return settingTime.UTC(), nil
}

// Get current time in Finnish time then convert to UTC
func GetCurrentTimeInUTC() (time.Time, error) {
	location, err := loadHelsinkiLocation()
	if err != nil {
		return time.Now(), err
	}
	now := time.Now().In(location)

	// Get year, month, and day components
	year, month, day := now.Date()
	currentTime := time.Date(year, month, day, now.Hour(), now.Minute(), 0, 0, location)
	// return as UTC
	return currentTime.UTC(), nil
}

// GetCurrentTimeInHelsinki returns the current time in Helsinki, Finland.
// It loads the Helsinki location and adjusts the current time to that timezone.
// The function returns the current time in Helsinki, the location object, and an error if any occurred during the location loading process.
func GetCurrentTimeInHelsinki() (time.Time, *time.Location, error) {
	location, err := loadHelsinkiLocation()
	if err != nil {
		return time.Now(), nil, err
	}
	now := time.Now().In(location)

	// Get year, month, and day components
	year, month, day := now.Date()
	currentTime := time.Date(year, month, day, now.Hour(), now.Minute(), 0, 0, location)
	// return as UTC
	return currentTime, location, nil
}

func loadHelsinkiLocation() (*time.Location, error) {
	location, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		return nil, fmt.Errorf("failed to get current location: %s", err.Error())
	}
	return location, nil
}
