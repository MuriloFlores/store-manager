package auth

import (
	"context"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type LoginUserUseCase struct {
	userRepository ports.UserRepository
	hasher         ports.PasswordHasher
	tokenManager   ports.TokenManager
	logger         ports.Logger
}

func NewLoginUserUseCase(userRepository ports.UserRepository, hasher ports.PasswordHasher, tokenManager ports.TokenManager, logger ports.Logger) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository: userRepository,
		hasher:         hasher,
		tokenManager:   tokenManager,
		logger:         logger,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	uc.logger.InfoLevel("Executing LoginUserUseCase with email: %s", map[string]interface{}{"email": email})

	user, err := uc.userRepository.FindByEmail(ctx, email)
	if err != nil {
		uc.logger.ErrorLevel("Failed to find user by email: %s, error: %v", err, map[string]interface{}{"email": email})
		return "", &domain.ErrInvalidCredentials{}
	}

	if !user.IsVerified() {
		uc.logger.InfoLevel("User email not verified: %s", map[string]interface{}{"email": email})
		return "", &domain.ErrEmailNotVerified{}
	}

	passwordMatch := uc.hasher.Compare(user.Password(), password)
	if !passwordMatch {
		uc.logger.ErrorLevel("Password mismatch for user: %s", err, map[string]interface{}{"email": email})
		return "", &domain.ErrInvalidCredentials{}
	}

	userIdentity := domain.Identity{
		UserID: user.ID(),
		Role:   user.Role(),
	}

	token, err := uc.tokenManager.Generate(&userIdentity)
	if err != nil {
		uc.logger.ErrorLevel("Failed to generate token for user: %s, error: %v", err, map[string]interface{}{"email": email})
		return "", fmt.Errorf("failed to generate token")
	}

	uc.logger.InfoLevel("Successfully logged in user: %s", map[string]interface{}{"email": email})

	return token, nil
}
