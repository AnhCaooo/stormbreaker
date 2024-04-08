package main

import (
	"log"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/AnhCaooo/stormbreaker/internal/routes"
	"github.com/gorilla/mux"
)

// todo: need handlings if request is not existing in routes
// todo: api docs
// todo: new endpoint?: next day, today price
// todo: improve logging as now it seems that there is no handlings for 404 NOT FOUND
func main() {
	// Initialize logger
	logger.InitLogger()

	// Initial new router
	r := mux.NewRouter()
	for _, endpoint := range routes.Endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	// Start server
	logger.Logger.Info("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
