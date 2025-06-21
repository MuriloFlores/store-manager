package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type UpdateUserUseCase struct {
	userRepository ports.UserRepository
	hasher         ports.PasswordHasher
}

type UpdateUserParams struct {
	Name *string
	Role *string
}

func NewUpdateUserUseCase(userRepository ports.UserRepository, hasher ports.PasswordHasher) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository: userRepository,
		hasher:         hasher,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, actor *domain.Identity, targetUserID string, params UpdateUserParams) (*domain.User, error) {
	targetUser, err := uc.userRepository.FindByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	if params.Name != nil {
		if actor.UserID != targetUserID && actor.Role != value_objects.Admin {
			return nil, &domain.ErrForbidden{Action: "attempt to update another user's name"}
		}

		if err = targetUser.ChangeName(*params.Name); err != nil {
			return nil, err
		}
	}

	if params.Role != nil {
		newRole := value_objects.Role(*params.Role)

		if err = uc.canChangeRole(actor, targetUser, newRole, ctx); err != nil {
			return nil, err
		}

		if err = targetUser.ChangeRole(newRole); err != nil {
			return nil, err
		}
	}

	if err = uc.userRepository.Update(ctx, targetUser); err != nil {
		return nil, err
	}

	return targetUser, nil
}

func (uc *UpdateUserUseCase) canChangeRole(actor *domain.Identity, targetUser *domain.User, newRole value_objects.Role, ctx context.Context) error {
	actorRole := actor.Role
	targetRole := targetUser.Role()

	switch actorRole {
	case value_objects.Admin:
		if actor.UserID == targetUser.ID() && newRole != value_objects.Admin {
			count, err := uc.userRepository.CountAdmins(ctx)
			if err != nil {
				return err
			}

			if count <= 1 {
				return &domain.ErrForbidden{Action: "attempt to change an admin"}
			}

		}

		return nil

	case value_objects.Manager:
		if targetRole == value_objects.Manager || targetRole == value_objects.Admin {
			return &domain.ErrForbidden{Action: "attempt to change an manager or admin"}
		}

		if targetRole == value_objects.Client {
			return &domain.ErrForbidden{Action: "attempt to change an client"}
		}

		if newRole == value_objects.Admin || newRole == value_objects.Manager {
			return &domain.ErrForbidden{Action: "attempt to change a user to admin or manager role"}
		}

		return nil
	}

	return &domain.ErrForbidden{Action: "attempt to update another user's role"}
}
