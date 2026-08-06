// Package middleware provides HTTP middleware for the Gateway's adapter/rest layer.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jeielsantos/aureon/modules/platform/internal/domain"
	"github.com/jeielsantos/aureon/modules/platform/internal/usecase"
)

type contextKey int

const apiKeyContextKey contextKey = iota

func RequireAPIKey(service *usecase.APIKeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret, ok := bearerToken(r)
			if !ok {
				http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
				return
			}

			key, err := service.Authenticate(r.Context(), secret)
			if err != nil {
				http.Error(w, "invalid or revoked api key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), apiKeyContextKey, key)))
		})
	}
}

func APIKeyFromContext(ctx context.Context) (domain.APIKey, bool) {
	key, ok := ctx.Value(apiKeyContextKey).(domain.APIKey)
	return key, ok
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}
