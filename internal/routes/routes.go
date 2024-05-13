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
		Path:    "/v1/market-price",
		Handler: api.PostMarketPrice,
		Method:  "POST",
	},
	{
		Path:    "/v1/market-price/today-tomorrow",
		Handler: api.GetTodayTomorrowPrice,
		Method:  "GET",
	},
}
