package helpers

import (
	"fmt"
	"time"
)

// return current date in Helsinki time with specific hour and minute.
// Use case examples: set expired time for caching
func SetTime(hour int, minute int) (time.Time, error) {
	// get timezone
	location, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		return time.Now(), fmt.Errorf("failed to get current location: %v", err)
	}

	// Get current time in Finnish time
	now := time.Now().In(location)

	// Get year, month, and day components
	year, month, day := now.Date()
	settingTime := time.Date(year, month, day, hour, minute, 0, 0, location)

	// Convert to UTC
	settingTimeInUTC := settingTime.UTC()

	return settingTimeInUTC, nil
}
