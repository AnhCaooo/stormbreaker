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
// todo: new endpoint?: next day, today price
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

	r.Use(middleware.LoggingMiddleware)
	// Start server
	logger.Logger.Info("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
