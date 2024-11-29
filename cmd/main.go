// AnhCao 2024
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AnhCaooo/go-goods/log"
	"github.com/AnhCaooo/stormbreaker/internal/api/handlers"
	"github.com/AnhCaooo/stormbreaker/internal/api/middleware"
	"github.com/AnhCaooo/stormbreaker/internal/api/routes"
	"github.com/AnhCaooo/stormbreaker/internal/config"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// todo: cache today-tomorrow price which means once the service starts, fetch and cache electric price
// and update the value when tomorrow price is available. Maybe have a service
// to listen and notify when the price is available. New service will also benefit for
// notifications service
func main() {
	ctx := context.Background()

	// todo: implement to accept a dynamic log level
	// Initialize logger
	logger := log.InitLogger(zapcore.InfoLevel)
	defer logger.Sync()

	configuration := &models.Config{}
	// Load and validate configuration
	err := config.LoadFile(configuration)
	if err != nil {
		logger.Fatal(constants.Server, zap.Error(err))
	}

	// todo: has own cache folder
	// Initialize resources: cache, database
	cache := models.NewCache()
	// Initialize database connection
	mongo := db.NewMongo(ctx, &configuration.Database, logger)
	mongoClient, err := mongo.EstablishConnection()
	if err != nil {
		logger.Fatal(constants.Server, zap.Error(err))
	}
	defer mongoClient.Disconnect(ctx)

	// Initialize Middleware
	middleware := middleware.NewMiddleware(logger, configuration)
	// Initialize Handler
	handler := handlers.NewHandler(logger, cache, mongo)
	// Initialize Endpoints pool
	endpoints := routes.InitializeEndpoints(handler)

	// Initial new router
	r := mux.NewRouter()

	middlewares := []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Authenticate,
	}
	for _, mw := range middlewares {
		r.Use(mw)
	}

	// Initialize pool of endpoints
	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handler.NotFound)

	// Start server
	logger.Info("Server started on", zap.String("port", configuration.Server.Port))
	http.ListenAndServe(fmt.Sprintf(":%s", configuration.Server.Port), r)
}
