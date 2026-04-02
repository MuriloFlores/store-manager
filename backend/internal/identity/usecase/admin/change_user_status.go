package admin

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/admin"
	"github.com/google/uuid"
)

type changeUserStatusUseCase struct {
	userRepo ports.UserRepository
}

func NewChangeUserStatusUseCase(userRepo ports.UserRepository) admin.ChangeUserStatusUseCase {
	return &changeUserStatusUseCase{userRepo: userRepo}
}

func (u *changeUserStatusUseCase) Execute(ctx context.Context, id string, active bool) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return entity.ErrUserNotFound
	}

	if active {
		user.Activate()
	} else {
		user.Deactivate()
	}

	return u.userRepo.Update(ctx, user)
}
