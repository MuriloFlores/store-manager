package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
)

type LoginUseCase struct {
	userRepo     ports.UserRepository
	tokenManager auth.TokenManager
	refreshRepo  ports.RefreshTokenRepository
	pepper       string
	expiresIn    time.Duration
}

func NewLogin(
	userRepo ports.UserRepository,
	tokenManager auth.TokenManager,
	refreshRepo ports.RefreshTokenRepository,
	pepper string,
	expiresIn time.Duration,
) auth.LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		refreshRepo:  refreshRepo,
		pepper:       pepper,
		expiresIn:    expiresIn,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input *dto.LoginRequest) (*dto.LoginResult, error) {
	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	if user == nil {
		return nil, entity.ErrInvalidCredentials
	}

	if ok := user.Password().Matches(input.Password, uc.pepper); !ok {
		return nil, entity.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := uc.tokenManager.GenerateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generating tokens: %w", err)
	}

	if err := uc.refreshRepo.SaveRefreshToken(ctx, user.ID(), refreshToken, uc.expiresIn); err != nil {
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	return &dto.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
