package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"go.uber.org/zap"

	"github.com/AnhCaooo/stormbreaker/internal/api/handlers"
	"github.com/AnhCaooo/stormbreaker/internal/api/middleware"
	"github.com/AnhCaooo/stormbreaker/internal/api/routes"
	"github.com/AnhCaooo/stormbreaker/internal/cache"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"github.com/gorilla/mux"
)

// API represents the main structure for the API server.
// It holds the configuration, context, logger, MongoDB connection,
// worker ID, HTTP server, and a wait group for managing goroutines.
type API struct {
	config   *models.Config
	ctx      context.Context
	logger   *zap.Logger
	mongo    *db.Mongo
	workerID int
	server   *http.Server
	wg       *sync.WaitGroup
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(ctx context.Context, logger *zap.Logger, config *models.Config, mongo *db.Mongo) *API {
	return &API{
		ctx:    ctx,
		config: config,
		logger: logger,
		mongo:  mongo,
	}
}

// Start initializes and starts the API server in a separate goroutine for a given worker.
// It sets up the server configuration, assigns the worker ID, and starts the server in a new goroutine.
// If the server encounters an error, it sends the error to the provided error channel.
func (a *API) Start(workerID int, errChan chan<- error, wg *sync.WaitGroup) {
	a.workerID = workerID
	a.wg = wg
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Server.Port),
		Handler: a.newMuxRouter(),
	}

	a.wg.Add(1)
	go func() {
		a.logger.Info(fmt.Sprintf("[worker_%d] Server starting...", a.workerID), zap.String("port", a.config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("[worker_%d] error in worker: %s", a.workerID, err.Error())
		}
	}()

}

// Shutdown the server gracefully
func (a *API) Stop() {
	defer a.wg.Done()
	a.logger.Info(fmt.Sprintf("[worker_%d] Stopping down HTTP server...", a.workerID), zap.String("port", a.config.Server.Port))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Fatal(fmt.Sprintf("[worker_%d] Server forced to shutdown", a.workerID), zap.Error(err))
	}

	a.logger.Info(fmt.Sprintf("[worker_%d] HTTP server stopped", a.workerID))
}

// todo: Proxy, CORS?
// newMuxRouter is responsible for all the top-level HTTP stuff that
// applies to all endpoints, like cache, database, CORS, auth middleware, and logging
func (a *API) newMuxRouter() *mux.Router {
	// Initialize cache
	cache := cache.NewCache(a.logger)
	// Initialize Middleware
	middleware := middleware.NewMiddleware(a.logger, a.config)
	// Initialize Handler
	apiHandler := handlers.NewHandler(a.logger, cache, a.mongo)
	// Initialize Endpoints pool
	endpoints := routes.InitializeEndpoints(apiHandler)

	r := mux.NewRouter()
	// Apply middlewares
	middlewares := []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Authenticate,
	}
	for _, mw := range middlewares {
		r.Use(mw)
	}

	// swagger endpoint for API documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// Apply endpoint handlers
	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	r.MethodNotAllowedHandler = http.HandlerFunc(apiHandler.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(apiHandler.NotFound)
	return r
}
