package middleware

import (
	"context"
	"net/http"

	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
	"github.com/gin-gonic/gin"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func RequireAuth(manager security.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		userClaims, err := manager.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid access token"})
			return
		}

		c.Set(UserClaimsKey, userClaims)

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), UserClaimsKey, userClaims))

		c.Next()
	}
}
