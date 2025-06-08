package middleware

import (
	"context"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

func AuthMiddleware(tokenManager ports.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				restErr := web_errors.NewUnauthorizedRequestError("Missing Authorization header")
				restErr.Send(w)
				return
			}

			identity, err := tokenManager.Validate(authHeader)
			if err != nil {
				restErr := web_errors.NewUnauthorizedRequestError("Invalid Authorization header")
				restErr.Send(w)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, identity.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, identity.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
