package auth

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type ConfirmUserEmailUseCase struct {
	userRepo  ports.UserRepository
	tokenRepo ports.ActionTokenRepository
	logger    ports.Logger
}

func NewConfirmUserEmailUseCase(userRepo ports.UserRepository, tokenRepo ports.ActionTokenRepository, logger ports.Logger) *ConfirmUserEmailUseCase {
	return &ConfirmUserEmailUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

func (uc *ConfirmUserEmailUseCase) Execute(ctx context.Context, tokenString string) error {
	uc.logger.InfoLevel("ConfirmUserEmailUseCase started", nil)

	actionToken, err := uc.tokenRepo.FindAndConsume(ctx, tokenString, domain.EmailConfirmation)
	if err != nil {
		uc.logger.ErrorLevel("ConfirmUserEmailUseCase error", err, map[string]interface{}{"token": tokenString})
		return err
	}

	newEmail := actionToken.Payload

	user, err := uc.userRepo.FindByID(ctx, actionToken.UserID)
	if err != nil {
		uc.logger.ErrorLevel("ConfirmUserEmailUseCase error", err, map[string]interface{}{"user_id": actionToken.UserID})
		return err
	}

	if err = user.ChangeEmail(newEmail); err != nil {
		uc.logger.ErrorLevel("ConfirmUserEmailUseCase error", err, map[string]interface{}{"user_id": user.ID, "new_email": newEmail})
		return err
	}

	if err = uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.ErrorLevel("ConfirmUserEmailUseCase error", err, map[string]interface{}{"user_id": user.ID, "new_email": newEmail})
		return err
	}

	uc.logger.InfoLevel("User email confirmed successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	return nil
}
