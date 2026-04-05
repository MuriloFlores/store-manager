package auth

import (
	"context"
	"fmt"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
	"github.com/google/uuid"
)

type changePasswordUseCase struct {
	userRepo ports.UserRepository
	logger   ports.Logger
	pepper   string
}

func NewChangePassword(
	userRepo ports.UserRepository,
	logger ports.Logger,
	pepper string,
) security.ChangePasswordUseCase {
	return &changePasswordUseCase{
		userRepo: userRepo,
		logger:   logger,
		pepper:   pepper,
	}
}

func (uc *changePasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	uc.logger.Debug("starting password change", "userID", userID)

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to find user for password change", err, "userID", userID)
		return fmt.Errorf("finding user by ID: %w", err)
	}
	if user == nil {
		uc.logger.Info("user not found for password change", "userID", userID)
		return entity.ErrUserNotFound
	}

	if !user.Password().Matches(oldPassword, uc.pepper) {
		uc.logger.Info("invalid old password provided", "userID", userID)
		return entity.ErrInvalidOldPassword
	}

	newPasswordVO, err := vo.NewPassword(newPassword, uc.pepper)
	if err != nil {
		uc.logger.Error("failed to create new password VO", err, "userID", userID)
		return fmt.Errorf("creating new password VO: %w", err)
	}

	user.ChangePassword(newPasswordVO)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("failed to update user password in repository", err, "userID", userID)
		return fmt.Errorf("updating user password: %w", err)
	}

	uc.logger.Info("password changed successfully", "userID", userID)
	return nil
}
