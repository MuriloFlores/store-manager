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
	pepper   string
}

func NewChangePassword(
	userRepo ports.UserRepository,
	pepper string,
) auth.ChangePasswordUseCase {
	return &changePasswordUseCase{
		userRepo: userRepo,
		pepper:   pepper,
	}
}

func (uc *changePasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding user by ID: %w", err)
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	newPasswordVO, err := vo.NewPassword(newPassword, uc.pepper)
	if err != nil {
		return fmt.Errorf("creating new password VO: %w", err)
	}

	if !user.Password().Matches(oldPassword, uc.pepper) {
		return entity.ErrInvalidOldPassword
	}

	user.ChangePassword(newPasswordVO)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("updating user password: %w", err)
	}

	return nil
}
