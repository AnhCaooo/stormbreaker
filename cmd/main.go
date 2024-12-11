// AnhCao 2024
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnhCaooo/go-goods/log"
	_ "github.com/AnhCaooo/stormbreaker/docs"
	"github.com/AnhCaooo/stormbreaker/internal/api/handlers"
	"github.com/AnhCaooo/stormbreaker/internal/api/middleware"
	"github.com/AnhCaooo/stormbreaker/internal/api/routes"
	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/config"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initializeRouter(handler *handlers.Handler, middleware *middleware.Middleware, endpoints []routes.Endpoint) *mux.Router {
	r := mux.NewRouter()

	// add swagger endpoint for API documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Apply middlewares
	middlewares := []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Authenticate,
	}
	for _, mw := range middlewares {
		r.Use(mw)
	}

	// Apply endpoint handlers
	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handler.NotFound)
	return r
}

func run(ctx context.Context, logger *zap.Logger, config *models.Config, r *mux.Router) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: r,
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", config.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	// Wait for termination signal
	select {
	case <-ctx.Done(): // Context cancellation
		logger.Warn("Context canceled")
	case <-stop: // OS signal received
		logger.Info("Termination signal received")
	}

	// Create a new context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

//	@title			Stormbreaker API (electric service)
//	@version		1.0.0
//	@description	Service for retrieving information about market electric price in Finland.

//	@contact.name	Anh Cao
//	@contact.email	anhcao4922@gmail.com

// @host		localhost:5001
// @BasePath	/
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

	// Initialize resources: cache, database
	cache := cache.NewCache(logger)
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

	// Initial new Mux router
	r := initializeRouter(handler, middleware, endpoints)
	// Start server
	run(ctx, logger, configuration, r)
}
