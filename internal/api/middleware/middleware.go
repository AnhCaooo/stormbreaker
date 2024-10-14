// AnhCao 2024
package middleware

import (
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Info("request received", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
