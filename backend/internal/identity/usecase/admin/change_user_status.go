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
	logger   ports.Logger
}

func NewChangeUserStatusUseCase(userRepo ports.UserRepository, logger ports.Logger) admin.ChangeUserStatusUseCase {
	return &changeUserStatusUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *changeUserStatusUseCase) Execute(ctx context.Context, id string, active bool) error {
	u.logger.Debug("starting change user status", "userID", id, "active", active)

	userID, err := uuid.Parse(id)
	if err != nil {
		u.logger.Error("failed to parse user ID", err, "id", id)
		return err
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to find user", err, "userID", userID)
		return err
	}

	if user == nil {
		u.logger.Info("user not found for status change", "userID", userID)
		return entity.ErrUserNotFound
	}

	if active {
		u.logger.Info("activating user", "userID", userID)
		user.Activate()
	} else {
		u.logger.Info("deactivating user", "userID", userID)
		user.Deactivate()
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		u.logger.Error("failed to update user status", err, "userID", userID)
		return err
	}

	u.logger.Info("user status updated successfully", "userID", userID, "active", active)
	return nil
}
