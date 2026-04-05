package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/redis/go-redis/v9"
)

type redisRateLimiter struct {
	client    *redis.Client
	limit     int
	window    time.Duration
	baseBlock time.Duration
}

func NewRedisLimiter(client *redis.Client, limit int, window, baseBlock time.Duration) ports.RateLimiterRepository {
	return &redisRateLimiter{
		client:    client,
		limit:     limit,
		window:    window,
		baseBlock: baseBlock,
	}
}

func (r *redisRateLimiter) Allow(ctx context.Context, key string) (bool, time.Duration, error) {
	now := time.Now().UnixNano()
	windowStart := now - r.window.Nanoseconds()

	zsetKey := fmt.Sprintf("rate_limit:%s", key)
	blockKey := fmt.Sprintf("rate_limit_block:%s", key)
	violationsKey := fmt.Sprintf("rate_limit_violation:%s", key)

	remainingBlock, err := r.client.TTL(ctx, zsetKey).Result()
	if err == nil && remainingBlock > 0 {
		return false, remainingBlock, nil
	}

	pipe := r.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, zsetKey, "0", fmt.Sprintf("%d", windowStart))
	pipe.ZCount(ctx, zsetKey, "-inf", "+inf")

	cmds, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, 0, err
	}

	count := cmds[1].(*redis.IntCmd).Val()

	if int(count) >= r.limit {
		violations, _ := r.client.Incr(ctx, violationsKey).Result()
		r.client.Expire(ctx, violationsKey, 24*time.Hour)

		exponent := math.Pow(2, float64(violations-1))
		blockDuration := time.Duration(float64(r.baseBlock) * exponent)

		r.client.Set(ctx, blockKey, "1", blockDuration)

		return false, blockDuration, nil
	}

	err = r.client.ZAdd(ctx, zsetKey, redis.Z{Score: float64(now), Member: now}).Err()
	if err != nil {
		return false, 0, err
	}

	r.client.Expire(ctx, zsetKey, r.window)

	return true, r.baseBlock, nil
}
