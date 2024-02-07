package main

import (
	"log"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/api"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Create a new router
	r := mux.NewRouter()

	// Define route for GET request to /data
	r.HandleFunc("/data", api.GetData).Methods("GET")

	// Start server
	logger.Logger.Info("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
