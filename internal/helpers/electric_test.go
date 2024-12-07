// AnhCao 2024
package helpers

import (
	"testing"
	"time"

	"github.com/AnhCaooo/stormbreaker/internal/models"
)

func TestFormatMarketPricePostReqParameters(t *testing.T) {
	tests := []struct {
		name           string
		requestPayload models.PriceRequest
		priceSettings  models.PriceSettings
		expectedUrl    string
		expectedErr    string
	}{
		{
			name: "valid request parameter",
			requestPayload: models.PriceRequest{
				StartDate:         "2024-06-05",
				EndDate:           "2024-06-05",
				Group:             "hour",
				CompareToLastYear: 0,
			},
			priceSettings: models.PriceSettings{
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedUrl: "https://oomi.fi/wp-json/spot-price/v1/get?starttime=2024-06-05&endtime=2024-06-05&margin=0.590000&group=hour&include_vat=1&compare_to_last_year=0",
			expectedErr: "",
		},
		{
			name: "invalid request parameter (invalid date range)",
			requestPayload: models.PriceRequest{
				StartDate:         "2024-06-07",
				EndDate:           "2024-06-05",
				Group:             "hour",
				CompareToLastYear: 0,
			},
			priceSettings: models.PriceSettings{
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedUrl: "",
			expectedErr: "start date cannot after end date",
		},
		{
			name: "invalid request parameter (invalid Group)",
			requestPayload: models.PriceRequest{
				StartDate:         "2024-06-05",
				EndDate:           "2024-06-05",
				Group:             "century",
				CompareToLastYear: 0,
			},
			priceSettings: models.PriceSettings{
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedUrl: "",
			expectedErr: "group should have valid value: 'hour', 'day', 'week', 'month', 'year'",
		},
		{
			name: "invalid request parameter (invalid CompareToLastYear)",
			requestPayload: models.PriceRequest{
				StartDate:         "2024-06-05",
				EndDate:           "2024-06-05",
				Group:             "hour",
				CompareToLastYear: 222,
			},
			priceSettings: models.PriceSettings{
				Marginal:    0.59,
				VatIncluded: true,
			},
			expectedUrl: "",
			expectedErr: "compareToLastYear needs to be value '0' or '1' only",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := FormatMarketPricePostReqParameters(&test.requestPayload, &test.priceSettings)
			if err != nil && err.Error() != test.expectedErr {
				t.Errorf("got %q, wanted %q", err.Error(), test.expectedErr)
			}

			if err == nil && result != test.expectedUrl {
				t.Errorf("got %q, wanted %q", result, test.expectedErr)
			}
		})

	}
}

func TestGetTodayAndTomorrowDateAsString(t *testing.T) {
	today, tomorrow := GetTodayAndTomorrowDateAsString()

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
