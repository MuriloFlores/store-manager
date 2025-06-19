package auth

import (
	"context"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type ConfirmPasswordResetUseCase struct {
	userRepo  ports.UserRepository
	tokenRepo ports.ActionTokenRepository
	hasher    ports.PasswordHasher
}

func NewConfirmPasswordResetUseCase(
	userRepo ports.UserRepository,
	tokenRepo ports.ActionTokenRepository,
	hasher ports.PasswordHasher,
) *ConfirmPasswordResetUseCase {
	return &ConfirmPasswordResetUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		hasher:    hasher,
	}
}

func (uc *ConfirmPasswordResetUseCase) Execute(ctx context.Context, tokenString string, newPassword string) error {
	actionToken, err := uc.tokenRepo.FindAndConsume(ctx, tokenString, domain.PasswordReset)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, actionToken.UserID)
	if err != nil {
		return err
	}

	newHashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("erro interno ao processar a nova senha: %w", err)
	}

	if err = user.SetPasswordHash(newHashedPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}
