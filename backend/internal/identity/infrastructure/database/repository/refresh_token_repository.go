package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type refreshTokenRepository struct {
	client *redis.Client
}

func NewRefreshTokenRepository(client *redis.Client) ports.RefreshTokenRepository {
	return &refreshTokenRepository{
		client: client,
	}
}

func (r *refreshTokenRepository) getKey(refreshToken string) string {
	return fmt.Sprintf("refresh_token:%s", refreshToken)
}

func (r *refreshTokenRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresIn time.Duration) error {
	return r.client.Set(ctx, r.getKey(refreshToken), userID.String(), expiresIn).Err()
}

func (r *refreshTokenRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	result, err := r.client.Get(ctx, r.getKey(refreshToken)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, entity.ErrSessionNotFound
		}

		return uuid.Nil, err
	}

	id, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	result, err := r.client.Del(ctx, r.getKey(refreshToken)).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		return entity.ErrSessionNotFound
	}

	return nil
}
