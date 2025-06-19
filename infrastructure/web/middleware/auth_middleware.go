package middleware

import (
	"context"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"net/http"
	"strings"
)

type contextKey string

const UserIdentityKey contextKey = "userIdentity"

func AuthMiddleware(tokenManager ports.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				restErr := web_errors.NewUnauthorizedRequestError("Missing Authorization header")
				restErr.Send(w)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				restErr := web_errors.NewUnauthorizedRequestError("invalid Authorization header format")
				restErr.Send(w)
				return
			}
			tokenString := headerParts[1]

			identity, err := tokenManager.Validate(tokenString)
			if err != nil {
				restErr := web_errors.NewUnauthorizedRequestError("expired or invalid token")
				restErr.Send(w)
				return
			}

			ctx := context.WithValue(r.Context(), UserIdentityKey, identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
