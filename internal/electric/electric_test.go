package electric

import (
	"math"
	"testing"
)

// TestIsValidFloat calls isValidFloat with a value, checking
// if value is not NaN and not infinite
func TestIsValidFloat(t *testing.T) {
	validValues := []float64{1.23, -45.67, 0.0, 1e-10, 1e20}
	invalidValues := []float64{math.NaN(), math.Inf(1), math.Inf(-1)}

	for _, value := range validValues {
		if !isValidFloat(value) {
			t.Errorf("Expected isValidFloat(%f) to be true, got false", value)
		}
	}

	for _, value := range invalidValues {
		if isValidFloat(value) {
			t.Errorf("Expected isValidFloat(%f) to be false, got true", value)
		}
	}
}

// TestIsValidDate calls isValidDate with a value, checking
// for a valid date value in string format.
func TestIsValidDate(t *testing.T) {
	validDates := []string{"2002-09-23", "2023-10-01"}
	invalidDates := []string{"string", "200222-09-21", "2002-13-43", ""}

	for _, dateStr := range validDates {
		if !isValidDate(dateStr) {
			t.Errorf("Expected isValidDate(%s) to be true, got false", dateStr)
		}
	}

	for _, dateStr := range invalidDates {
		if isValidDate(dateStr) {
			t.Errorf("Expected isValidDate(%s) to be false, got true", dateStr)
		}
	}
}

// TestIsValidGroup calls isValidGroup with a value, checking
// for a valid group: 'hour', 'day', 'week', 'month', 'year'.
func TestIsValidGroup(t *testing.T) {
	validGroups := []string{"hour", "day", "week", "month", "year"}
	invalidGroups := []string{"invalid-value", "hours", "weeks"}

	for _, group := range validGroups {
		if !isValidGroup(group) {
			t.Errorf("Expected isValidGroup(%s) to be true, got false", group)
		}
	}

	for _, group := range invalidGroups {
		if isValidGroup(group) {
			t.Errorf("Expected isValidGroup(%s) to be false, got true", group)
		}
	}
}

// TestIsValidInt calls isValidInt with a value, checking
// for a valid int value: "0" or "1".
func TestIsValidInt(t *testing.T) {
	validValues := []int32{0, 1}
	invalidValues := []int32{3, 43, 324}

	for _, intVal := range validValues {
		if !isValidInt(intVal) {
			t.Errorf("Expected isValidInt(%d) to be true, got false", intVal)
		}
	}

	for _, intVal := range invalidValues {
		if isValidInt(intVal) {
			t.Errorf("Expected isValidInt(%d) to be false, got true", intVal)
		}
	}
}
