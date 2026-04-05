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
	tokenManager security.TokenManager
	refreshRepo  ports.RefreshTokenRepository
	logger       ports.Logger
	pepper       string
	baseDuration time.Duration
	threshold    int
	expiresIn    time.Duration
}

func NewLogin(
	userRepo ports.UserRepository,
	tokenManager security.TokenManager,
	refreshRepo ports.RefreshTokenRepository,
	logger ports.Logger,
	pepper string,
	baseDuration time.Duration,
	threshold int,
	expiresIn time.Duration,
) security.LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		refreshRepo:  refreshRepo,
		logger:       logger,
		pepper:       pepper,
		baseDuration: baseDuration,
		threshold:    threshold,
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

	if user.IsLocked(time.Now()) {
		uc.logger.Info("user is locked", "email", input.Email)
		return nil, entity.ErrUserBlocked
	}

	if ok := user.Password().Matches(input.Password, uc.pepper); !ok {
		uc.logger.Info("login failed: invalid password", "userID", user.ID())

		user.RecordFailedLogin(uc.threshold, uc.baseDuration, time.Now())
		err := uc.userRepo.Update(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("error updating failed attempts: %w", err)
		}

		return nil, entity.ErrInvalidCredentials
	}

	user.ResetFailedAttempts()

	accessToken, refreshToken, err := uc.tokenManager.GenerateTokens(ctx, user)
	if err != nil {
		uc.logger.Error("failed to generate auth tokens", err, "userID", user.ID())
		return nil, fmt.Errorf("generating tokens: %w", err)
	}

	if err := uc.refreshRepo.SaveRefreshToken(ctx, user.ID(), refreshToken, uc.expiresIn); err != nil {
		uc.logger.Error("failed to save refresh token", err, "userID", user.ID())
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error reset failed attempts: %w", err)
	}

	uc.logger.Info("user logged in successfully", "userID", user.ID())
	return &dto.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
