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
	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/config"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
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

	// Initialize logger
	logger := log.InitLogger(zapcore.InfoLevel)
	defer logger.Sync()

	// Read configuration file
	err := config.ReadFile(&config.Config)
	if err != nil {
		logger.Fatal(constants.Server, zap.Error(err))
	}

	// Initialize cache
	cache.NewCache()

	// Initialize database connection
	mongo := db.Init(ctx, &config.Config.Database, logger, nil)
	mongoClient, err := mongo.EstablishConnection()
	if err != nil {
		logger.Fatal(constants.Server, zap.Error(err))
	}
	defer mongoClient.Disconnect(ctx)

	// Initial new router
	r := mux.NewRouter()
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Authenticate)

	for _, endpoint := range routes.Endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	// Start server
	logger.Info("Server started on", zap.String("port", config.Config.Server.Port))
	http.ListenAndServe(fmt.Sprintf(":%s", config.Config.Server.Port), r)
}
