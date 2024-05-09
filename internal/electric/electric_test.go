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

// TestIsValidDateRange calls isValidDate with a value, checking
// for a valid date value in string format.
func TestIsValidDateRange(t *testing.T) {
	tests := []struct {
		startDate string
		endDate   string
		expected  bool
	}{
		{"2022-01-01", "2022-01-02", true},
		{"2022-05-02", "2022-05-01", false},
		{"2022-05-09", "2022-05-09", true},
		{"2022-13-02", "2022-03-01", false}, // invalid start date
		{"2022-12-02", "2022-19-01", false}, // invalid end date
		{"2022-01-01", "invalid", false},
		{"invalid", "2022-01-02", false},
	}

	for _, test := range tests {
		actual, err := isValidDateRange(test.startDate, test.endDate)
		if err != nil && actual != test.expected {
			t.Errorf("error while validating input (%s, %s): got %v", test.startDate, test.endDate, err)
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
