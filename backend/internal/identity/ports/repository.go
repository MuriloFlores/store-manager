package ports

import (
	"context"
	"time"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error)
	GetUsersInfo(ctx context.Context, roles []vo.Role, pagination _common.Pagination) (*_common.PaginatedResult[*entity.User], error)
	Update(ctx context.Context, user *entity.User) error
}

type RefreshTokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresIn time.Duration) error
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type OTPRepository interface {
	SaveOTP(ctx context.Context, email vo.Email, otp vo.OTP, expiresIn time.Duration) error
	GetOTP(ctx context.Context, email vo.Email) (vo.OTP, error)
	DeleteOTP(ctx context.Context, email vo.Email) error
}
