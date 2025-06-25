package ports

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit time.Duration) (bool, error)
}
