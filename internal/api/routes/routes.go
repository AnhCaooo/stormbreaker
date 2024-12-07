// AnhCao 2024
package routes

import (
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/api/handlers"
)

// Endpoint is the presentation of object which contains values for routing
type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

// InitializeEndpoints creates a pool of Endpoints
func InitializeEndpoints(handler *handlers.Handler) []Endpoint {
	return []Endpoint{
		{
			Path:    "/v1/ping",
			Handler: handler.Ping,
			Method:  "GET",
		},
		{
			Path:    "/v1/market-price",
			Handler: handler.PostMarketPrice,
			Method:  "POST",
		},
		{
			Path:    "/v1/market-price/today-tomorrow",
			Handler: handler.GetTodayTomorrowPrice,
			Method:  "GET",
		},
		{
			Path:    "/v1/price-settings",
			Handler: handler.GetPriceSettings,
			Method:  "GET",
		},
		{
			Path:    "/v1/price-settings",
			Handler: handler.CreatePriceSettings,
			Method:  "POST",
		},
		{
			Path:    "/v1/price-settings",
			Handler: handler.PatchPriceSettings,
			Method:  "PATCH",
		},
		{
			Path:    "/v1/price-settings",
			Handler: handler.DeletePriceSettings,
			Method:  "DELETE",
		},
		// ? /v1/market-price/usage-situation - use AI to analyze from which time user can use normally, or just fixed limit?
	}
}
