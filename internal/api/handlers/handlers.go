// AnhCao 2024
package handlers

import (
	"fmt"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	cache  *models.Cache
	mongo  *db.Mongo
}

func NewHandler(logger *zap.Logger, cache *models.Cache, mongo *db.Mongo) *Handler {
	return &Handler{
		logger: logger,
		cache:  cache,
		mongo:  mongo,
	}
}

// return response when request url is not found
func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.logger.Info("undefined endpoint", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("404 - Not found"))
}

// return response when request method is not allowed
func (h Handler) NotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	h.logger.Info("method not allowed", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("405 - Method not allowed"))
}

// Ping the connection to the server
func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func (h Handler) Logger(next http.Handler, w http.ResponseWriter, r *http.Request) {
	h.logger.Info("request received", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	next.ServeHTTP(w, r)
}
