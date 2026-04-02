package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/redis/go-redis/v9"
)

type otpRepository struct {
	client *redis.Client
}

func NewOtpRepository(client *redis.Client) ports.OTPRepository {
	return &otpRepository{client: client}
}

func (o *otpRepository) getKey(email vo.Email) string {
	return fmt.Sprintf("otp:%s", email.String())
}

func (o *otpRepository) SaveOTP(ctx context.Context, email vo.Email, otp vo.OTP, expiresIn time.Duration) error {
	return o.client.Set(ctx, o.getKey(email), otp.String(), expiresIn).Err()
}

func (o *otpRepository) GetOTP(ctx context.Context, email vo.Email) (vo.OTP, error) {
	result, err := o.client.Get(ctx, o.getKey(email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", entity.ErrOTPNotFound
		}

		return "", err
	}

	otpVo, err := vo.NewOTP(result)
	if err != nil {
		return "", err
	}

	return otpVo, nil
}

func (o *otpRepository) DeleteOTP(ctx context.Context, email vo.Email) error {
	result, err := o.client.Del(ctx, o.getKey(email)).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		return entity.ErrOTPNotFound
	}

	return nil
}
