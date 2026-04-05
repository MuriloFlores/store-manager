package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
)

type rotateRefreshTokenUseCase struct {
	userRepo     ports.UserRepository
	refreshRepo  ports.RefreshTokenRepository
	tokenManager security.TokenManager
	logger       ports.Logger
	expiresIn    time.Duration
}

func NewRotateRefreshTokenUseCase(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	tokenManger security.TokenManager,
	logger ports.Logger,
	expiresIn time.Duration,
) security.RotateRefreshTokenUseCase {
	return &rotateRefreshTokenUseCase{
		userRepo:     userRepo,
		refreshRepo:  refreshRepo,
		tokenManager: tokenManger,
		logger:       logger,
		expiresIn:    expiresIn,
	}
}

func (uc *rotateRefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*dto.LoginResult, error) {
	uc.logger.Debug("starting refresh token rotation")

	userID, err := uc.refreshRepo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		uc.logger.Info("failed to get user ID from refresh token (invalid or expired)", "error", err)
		return nil, fmt.Errorf("getting user ID from refresh token: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to find user during token rotation", err, "userID", userID)
		return nil, fmt.Errorf("finding user by ID: %w", err)
	}

	if user == nil {
		uc.logger.Info("user not found during token rotation", "userID", userID)
		return nil, entity.ErrUserNotFound
	}

	if !user.IsActive() {
		uc.logger.Info("token rotation failed: user is deactivated", "userID", userID)
		return nil, entity.ErrUserIsDeactivated
	}

	if err := uc.refreshRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		uc.logger.Error("failed to delete old refresh token", err, "userID", userID)
		return nil, fmt.Errorf("deleting old refresh token: %w", err)
	}

	accessToken, refreshToken, err := uc.tokenManager.GenerateTokens(ctx, user)
	if err != nil {
		uc.logger.Error("failed to generate new tokens", err, "userID", userID)
		return nil, fmt.Errorf("generating new tokens: %w", err)
	}

	if err := uc.refreshRepo.SaveRefreshToken(ctx, user.ID(), refreshToken, uc.expiresIn); err != nil {
		uc.logger.Error("failed to save new refresh token", err, "userID", userID)
		return nil, fmt.Errorf("saving new refresh token: %w", err)
	}

	uc.logger.Info("tokens rotated successfully", "userID", userID)
	return &dto.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
