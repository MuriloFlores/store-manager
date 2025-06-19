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
}

func NewLoginUserUseCase(userRepository ports.UserRepository, hasher ports.PasswordHasher, tokenManager ports.TokenManager) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository: userRepository,
		hasher:         hasher,
		tokenManager:   tokenManager,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return "", &domain.ErrInvalidCredentials{}
	}

	if !user.IsVerified() {
		return "", &domain.ErrEmailNotVerified{}
	}

	passwordMatch := uc.hasher.Compare(user.Password(), password)
	if !passwordMatch {
		return "", &domain.ErrInvalidCredentials{}
	}

	userIdentity := domain.Identity{
		UserID: user.ID(),
		Role:   user.Role(),
	}

	token, err := uc.tokenManager.Generate(&userIdentity)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}
