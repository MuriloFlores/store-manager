package auth

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/security"
)

type logoutUseCase struct {
	refreshRepo ports.RefreshTokenRepository
	logger      ports.Logger
}

func NewLogoutUseCase(refreshRepo ports.RefreshTokenRepository, logger ports.Logger) security.LogoutUseCase {
	return &logoutUseCase{
		refreshRepo: refreshRepo,
		logger:      logger,
	}
}

func (uc *logoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	uc.logger.Debug("starting logout process")

	if err := uc.refreshRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		uc.logger.Error("failed to delete refresh token during logout", err)
		return err
	}

	uc.logger.Info("user logged out successfully")
	return nil
}
