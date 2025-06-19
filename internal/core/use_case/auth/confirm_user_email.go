package auth

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type ConfirmUserEmailUseCase struct {
	userRepo  ports.UserRepository
	tokenRepo ports.ActionTokenRepository
}

func NewConfirmUserEmailUseCase(userRepo ports.UserRepository, tokenRepo ports.ActionTokenRepository) *ConfirmUserEmailUseCase {
	return &ConfirmUserEmailUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (uc *ConfirmUserEmailUseCase) Execute(ctx context.Context, tokenString string) error {
	actionToken, err := uc.tokenRepo.FindAndConsume(ctx, tokenString, domain.EmailConfirmation)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, actionToken.UserID)
	if err != nil {
		return err
	}

	user.MarkAsVerified()

	return uc.userRepo.Update(ctx, user)
}
