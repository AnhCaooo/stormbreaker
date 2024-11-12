// AnhCao 2024
package middleware

import (
	"net/http"
	"strings"

	"github.com/AnhCaooo/go-goods/auth"
	"github.com/AnhCaooo/stormbreaker/internal/config"
	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

// log the coming request to the server
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Info("request received", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

// read the token from request and do verify the access token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusForbidden)
			logger.Logger.Info("permission Denied: No token provided")
			w.Write([]byte("403 - Forbidden"))
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		err := auth.VerifyToken(tokenString, config.Config.Supabase.Auth.JwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Logger.Error("unauthorized request", zap.Error(err))
			w.Write([]byte("401 - Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
