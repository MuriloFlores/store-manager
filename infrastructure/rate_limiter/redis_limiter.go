package rate_limiter

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) ports.RateLimiter {
	return &RedisRateLimiter{
		client: client,
	}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit time.Duration) (bool, error) {
	wasSet, err := r.client.SetNX(ctx, key, true, limit).Result()
	if err != nil {
		return false, err
	}

	return wasSet, nil
}
