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
	logger       ports.Logger
	pepper       string
	expiresIn    time.Duration
}

func NewLogin(
	userRepo ports.UserRepository,
	tokenManager auth.TokenManager,
	refreshRepo ports.RefreshTokenRepository,
	logger ports.Logger,
	pepper string,
	expiresIn time.Duration,
) auth.LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		refreshRepo:  refreshRepo,
		logger:       logger,
		pepper:       pepper,
		expiresIn:    expiresIn,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input *dto.LoginRequest) (*dto.LoginResult, error) {
	uc.logger.Debug("starting login process", "email", input.Email)

	email, err := vo.NewEmail(input.Email)
	if err != nil {
		uc.logger.Info("invalid email format in login", "email", input.Email)
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		uc.logger.Error("error finding user during login", err, "email", input.Email)
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	if user == nil {
		uc.logger.Info("login failed: user not found", "email", input.Email)
		return nil, entity.ErrInvalidCredentials
	}

	if ok := user.Password().Matches(input.Password, uc.pepper); !ok {
		uc.logger.Info("login failed: invalid password", "userID", user.ID())
		return nil, entity.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := uc.tokenManager.GenerateTokens(ctx, user)
	if err != nil {
		uc.logger.Error("failed to generate auth tokens", err, "userID", user.ID())
		return nil, fmt.Errorf("generating tokens: %w", err)
	}

	if err := uc.refreshRepo.SaveRefreshToken(ctx, user.ID(), refreshToken, uc.expiresIn); err != nil {
		uc.logger.Error("failed to save refresh token", err, "userID", user.ID())
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	uc.logger.Info("user logged in successfully", "userID", user.ID())
	return &dto.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
