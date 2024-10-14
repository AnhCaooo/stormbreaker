// AnhCao 2024
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AnhCaooo/stormbreaker/internal/api/handlers"
	"github.com/AnhCaooo/stormbreaker/internal/api/middleware"
	"github.com/AnhCaooo/stormbreaker/internal/api/routes"
	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/config"
	title "github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// todo: cache today-tomorrow price which means once the service starts, fetch and cache electric price
// and update the value when tomorrow price is available. Maybe have a service
// to listen and notify when the price is available. New service will also benefit for
// notifications service
func main() {
	ctx := context.Background()

	// Initialize logger
	logger.InitLogger()

	// Read configuration file
	err := config.ReadFile(&config.Config)
	if err != nil {
		logger.Logger.Error(title.Server, zap.Error(err))
		os.Exit(1)
	}

	// Initialize cache
	cache.NewCache()

	// Initialize database connection
	mongo, err := db.Init(ctx, config.Config.Database)
	if err != nil {
		logger.Logger.Error(title.Server, zap.Error(err))
		os.Exit(1)
	}
	defer mongo.Disconnect(ctx)

	// Initial new router
	r := mux.NewRouter()
	for _, endpoint := range routes.Endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	// Middleware
	r.Use(middleware.Logger)

	// Start server
	logger.Logger.Info("Server started on", zap.String("port", config.Config.Server.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Config.Server.Port), r))
}
