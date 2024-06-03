package api

import (
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

// return response when request url is not found
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	logger.Logger.Info("undefined endpoint", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("404 - Not found"))
}

// return response when request method is not allowed
func NotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	logger.Logger.Info("method not allowed", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("405 - Method not allowed"))
}
