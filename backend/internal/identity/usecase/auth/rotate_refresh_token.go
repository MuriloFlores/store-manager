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
	tokenManager auth.TokenManager
	expiresIn    time.Duration
}

func NewRotateRefreshTokenUseCase(
	userRepo ports.UserRepository,
	refreshRepo ports.RefreshTokenRepository,
	tokenManger auth.TokenManager,
	expiresIn time.Duration,
) auth.RotateRefreshTokenUseCase {
	return &rotateRefreshTokenUseCase{
		userRepo:     userRepo,
		refreshRepo:  refreshRepo,
		tokenManager: tokenManger,
		expiresIn:    expiresIn,
	}
}

func (uc *rotateRefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*dto.LoginResult, error) {
	userID, err := uc.refreshRepo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("getting user ID from refresh token: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("finding user by ID: %w", err)
	}

	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	if !user.IsActive() {
		return nil, entity.ErrUserIsDeactivated
	}

	if err := uc.refreshRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("deleting old refresh token: %w", err)
	}

	accessToken, refreshToken, err := uc.tokenManager.GenerateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generating new tokens: %w", err)
	}

	if err := uc.refreshRepo.SaveRefreshToken(ctx, user.ID(), refreshToken, uc.expiresIn); err != nil {
		return nil, fmt.Errorf("saving new refresh token: %w", err)
	}

	return &dto.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
