package auth

import (
	"context"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type ChangePasswordUseCase struct {
	userRepo     repositories.UserRepository
	hasher       ports.PasswordHasher
	taskEnqueuer ports.TaskEnqueuer
}

func NewChangePasswordUseCase(userRepo repositories.UserRepository, hasher ports.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, actor *domain.Identity, oldPassword, newPassword string) error {
	user, err := uc.userRepo.FindByID(ctx, actor.UserID)
	if err != nil {
		return &domain.ErrInvalidCredentials{}
	}

	if !uc.hasher.Compare(user.Password(), oldPassword) {
		return &domain.ErrInvalidCredentials{}
	}

	newHashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("internal server error in password processing: %w", err)
	}

	if err = user.SetPasswordHash(newHashedPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}
