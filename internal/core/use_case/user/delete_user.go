package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type DeleteUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewDeleteUserUseCase(userRepository repositories.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepository: userRepository}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, actor *domain.Identity, targetUserID string) error {
	targetUser, err := uc.userRepository.FindByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	isOwner := actor.UserID == targetUser.ID()
	isAdmin := actor.Role == value_objects.Admin

	if !isOwner && !isAdmin {
		return &domain.ErrForbidden{Action: "attempt to delete a third party user"}
	}

	if targetUser.Role() == value_objects.Admin {
		count, err := uc.userRepository.CountAdmins(ctx)
		if err != nil {
			return err
		}

		if count <= 1 {
			return &domain.ErrForbidden{Action: "attempt to delete admin user without permission"}
		}
	}

	return uc.userRepository.Delete(ctx, targetUser.ID())
}
