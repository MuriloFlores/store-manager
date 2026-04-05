package middleware

import (
	"fmt"
	"net/http"

	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/gin-gonic/gin"
)

func RateLimit(limiter ports.RateLimiterRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		allowed, retryAfter, err := limiter.Allow(c.Request.Context(), key)
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":               "too many requests, please try again later",
				"retry_after_seconds": retryAfter.Seconds(),
			})

			return
		}

		c.Next()
	}
}
