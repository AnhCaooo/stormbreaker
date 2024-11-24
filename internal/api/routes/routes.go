// AnhCao 2024
package routes

import (
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/api"
)

// Endpoint is the presentation of object which contains values for routing
type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

var Endpoints = []Endpoint{
	{
		Path:    "/v1/ping",
		Handler: api.Ping,
		Method:  "GET",
	},
	{
		Path:    "/v1/market-price",
		Handler: api.PostMarketPrice,
		Method:  "POST",
	},
	{
		Path:    "/v1/market-price/today-tomorrow",
		Handler: api.GetTodayTomorrowPrice,
		Method:  "GET",
	},
	{
		Path:    "v1/price-settings",
		Handler: nil,
		Method:  "GET",
	},
	{
		Path:    "v1/price-settings",
		Handler: nil,
		Method:  "POST",
	},
	{
		Path:    "v1/price-settings",
		Handler: nil,
		Method:  "PATCH",
	},
	{
		Path:    "v1/price-settings",
		Handler: nil,
		Method:  "DELETE",
	},
	// ? /v1/market-price/usage-situation - use AI to analyze from which time user can use normally, or just fixed limit?
}
