package main

import (
	"log"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/api"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/middleware"
	"github.com/AnhCaooo/stormbreaker/internal/routes"
	"github.com/gorilla/mux"
)

// todo: api docs
// todo: cache today-tomorrow price which means once the service starts, fetch and cache electric price
// and update the value when tomorrow price is available. Maybe have a service
// to listen and notify when the price is available. New service will also benefit for
// notifications service
func main() {
	// Initialize logger
	logger.InitLogger()

	// Initial new router
	r := mux.NewRouter()
	for _, endpoint := range routes.Endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}
	r.MethodNotAllowedHandler = http.HandlerFunc(api.NotAllowedHandler)
	r.NotFoundHandler = http.HandlerFunc(api.NotFoundHandler)

	// Middleware
	r.Use(middleware.LoggingMiddleware)

	// Start server
	logger.Logger.Info("Server started on :5001")
	log.Fatal(http.ListenAndServe(":5001", r))
}
