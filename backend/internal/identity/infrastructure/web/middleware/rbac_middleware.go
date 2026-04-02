package middleware

import (
	"net/http"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/web/helper"
	"github.com/gin-gonic/gin"
)

func VerifyRole(allowedRoles ...vo.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := helper.ExtractUserClaims(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "auth context missing or invalid token"})
			return
		}

		allowedMap := make(map[vo.Role]bool)
		for _, role := range allowedRoles {
			allowedMap[role] = true
		}

		hasPermission := false

		for _, userRole := range claims.Roles {
			if allowedMap[userRole] {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden: insufficient permissions"})
			return
		}

		c.Next()
	}
}
