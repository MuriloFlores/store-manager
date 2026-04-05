package security

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
)

type RotateRefreshTokenUseCase interface {
	Execute(ctx context.Context, refreshToken string) (*dto.LoginResult, error)
}
