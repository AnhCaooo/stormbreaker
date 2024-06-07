package helpers

import (
	"testing"
	"time"
)

func TestGetTodayAndTomorrowDateAsString(t *testing.T) {
	today, tomorrow := getTodayAndTomorrowDateAsString()

	// Get current date and time
	now := time.Now()

	// Calculate expected today and tomorrow dates
	expectedToday := now.Truncate(24 * time.Hour).Format(DATE_FORMAT)
	expectedTomorrow := now.Truncate(24*time.Hour).AddDate(0, 0, 1).Format(DATE_FORMAT)

	if today != expectedToday {
		t.Errorf("Expected today: %s, but got: %s", expectedToday, today)
	}

	if tomorrow != expectedTomorrow {
		t.Errorf("Expected tomorrow: %s, but got: %s", expectedTomorrow, tomorrow)
	}
}
