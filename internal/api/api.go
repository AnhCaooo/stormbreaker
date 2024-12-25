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

func (a *API) Start(workerID int, errChan chan<- error, wg *sync.WaitGroup) {
	a.workerID = workerID
	a.wg = wg
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Server.Port),
		Handler: a.newMuxRouter(),
	}

	go func() {
		a.logger.Info("Server starting", zap.String("port", a.config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("error in worker %d: %s", 1, err.Error())
		}
	}()

}

// Shutdown the server gracefully
func (a *API) Stop() {
	defer a.wg.Done()
	a.logger.Info("Stopping down HTTP server in worker", zap.Int("worker_id", a.workerID))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	a.logger.Info("HTTP server stopped", zap.Int("worker_id", a.workerID))

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
