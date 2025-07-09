package middleware

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"net/http"
	"strings"
)

func TryAuthMiddleware(tokenManager ports.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := headerParts[1]

			identity, err := tokenManager.Validate(tokenString)
			if err == nil && identity != nil {
				ctx := context.WithValue(r.Context(), UserIdentityKey, identity)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
