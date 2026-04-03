package admin

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/admin"
	"github.com/google/uuid"
)

type changeUserRoleUseCase struct {
	userRepo ports.UserRepository
	logger   ports.Logger
}

func NewChangeUserRoleUseCase(userRepo ports.UserRepository, logger ports.Logger) admin.ChangeUserRoleUseCase {
	return &changeUserRoleUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *changeUserRoleUseCase) Execute(ctx context.Context, id string, roles []string) error {
	u.logger.Debug("starting change user roles", "userID", id, "newRoles", roles)

	userID, err := uuid.Parse(id)
	if err != nil {
		u.logger.Error("failed to parse user ID", err, "id", id)
		return err
	}

	rolesVo := make([]vo.Role, 0, len(roles))
	for _, r := range roles {
		validRole, err := vo.NewRole(r)
		if err != nil {
			u.logger.Error("invalid role provided", err, "role", r)
			return err
		}

		rolesVo = append(rolesVo, validRole)
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to find user", err, "userID", userID)
		return err
	}

	if user == nil {
		u.logger.Info("user not found for role change", "userID", userID)
		return entity.ErrUserNotFound
	}

	u.logger.Info("updating user roles", "userID", userID, "oldRoles", user.Roles(), "newRoles", rolesVo)
	user.ReplaceRoles(rolesVo)

	if err := u.userRepo.Update(ctx, user); err != nil {
		u.logger.Error("failed to update user roles", err, "userID", userID)
		return err
	}

	u.logger.Info("user roles updated successfully", "userID", userID)
	return nil
}
