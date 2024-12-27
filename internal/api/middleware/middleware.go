// AnhCao 2024
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AnhCaooo/go-goods/auth"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"go.uber.org/zap"
)

type Middleware struct {
	logger   *zap.Logger
	config   *models.Config
	workerID int
}

func NewMiddleware(logger *zap.Logger, config *models.Config, workerID int) *Middleware {
	return &Middleware{
		logger:   logger,
		config:   config,
		workerID: workerID,
	}
}

// log the coming request to the server
func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}
		m.logger.Info(fmt.Sprintf("[worker_%d] request received", m.workerID), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

// read the token from request and do verify the access token
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			m.logger.Error(fmt.Sprintf("[worker_%d] permission Denied: No authentication provided in header", m.workerID))
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := auth.VerifyToken(tokenString, m.config.Supabase.Auth.JwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			m.logger.Error(fmt.Sprintf("[worker_%d] unauthorized request", m.workerID), zap.Error(err))
			w.Write([]byte("401 - Unauthorized"))
			return
		}

		// due to 'Supabase' authentication, it stores userId via "sub" field
		userID, err := auth.ExtractValueFromTokenClaim(token, "sub")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			m.logger.Error(fmt.Sprintf("[worker_%d] unauthorized request", m.workerID), zap.Error(err))
			w.Write([]byte("401 - Unauthorized"))
			return
		}

		// Add userID to the context
		ctx := context.WithValue(r.Context(), constants.UserIdKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
