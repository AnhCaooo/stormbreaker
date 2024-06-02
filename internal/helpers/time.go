package helpers

import "time"

// return current date with specific hour and minute.
// Use case examples: set expired time for caching
func SetTime(hour int, minute int) time.Time {
	// Get current time
	now := time.Now()

	// Get year, month, and day components
	year, month, day := now.Date()

	return time.Date(year, month, day, hour, minute, 0, 0, now.Location())
}
