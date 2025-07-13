package auth

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type ConfirmAccountUserUseCase struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.ActionTokenRepository
	logger    ports.Logger
}

func NewConfirmAccountUserUseCase(
	UserRepo repositories.UserRepository,
	tokenRepo repositories.ActionTokenRepository,
	logger ports.Logger,
) *ConfirmAccountUserUseCase {
	return &ConfirmAccountUserUseCase{
		userRepo:  UserRepo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

func (uc *ConfirmAccountUserUseCase) Execute(ctx context.Context, tokenString string) error {
	uc.logger.InfoLevel("ConfirmAccountUserUseCase started", nil)

	actionToken, err := uc.tokenRepo.FindAndConsume(ctx, tokenString, domain.AccountVerification)
	if err != nil {
		uc.logger.ErrorLevel("ConfirmAccountUserUseCase error", err, map[string]interface{}{"token": tokenString})
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, actionToken.UserID)
	if err != nil {
		uc.logger.ErrorLevel("ConfirmAccountUserUseCase error", err, map[string]interface{}{"user_id": actionToken.UserID})
		return err
	}

	user.MarkAsVerified()

	uc.logger.InfoLevel("User account confirmed successfully", map[string]interface{}{"user_id": user.ID(), "email": user.Email()})

	return uc.userRepo.Update(ctx, user)
}
