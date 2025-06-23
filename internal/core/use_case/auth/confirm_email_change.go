package auth

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type ConfirmEmailChangeUseCase struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.ActionTokenRepository
}

func NewConfirmEmailChangeUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.ActionTokenRepository,
) *ConfirmEmailChangeUseCase {
	return &ConfirmEmailChangeUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (uc *ConfirmEmailChangeUseCase) Execute(ctx context.Context, tokenString string) error {
	actionToken, err := uc.tokenRepo.FindAndConsume(ctx, tokenString, domain.EmailConfirmation)
	if err != nil {
		return err
	}

	newEmail := actionToken.Payload

	user, err := uc.userRepo.FindByID(ctx, actionToken.UserID)
	if err != nil {
		return err
	}

	if err = user.ChangeEmail(newEmail); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}
