package auth

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
)

type logoutUseCase struct {
	refreshRepo ports.RefreshTokenRepository
}

func NewLogoutUseCase(refreshRepo ports.RefreshTokenRepository) auth.LogoutUseCase {
	return &logoutUseCase{
		refreshRepo: refreshRepo,
	}
}

func (uc *logoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	return uc.refreshRepo.DeleteRefreshToken(ctx, refreshToken)
}
