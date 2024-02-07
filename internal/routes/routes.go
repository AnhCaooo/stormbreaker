package routes

import (
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/api"
)

type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

var Endpoints = []Endpoint{
	{
		Path:    "/v1/market-price",
		Handler: api.GetMarketPrice,
		Method:  "GET",
	},
}
