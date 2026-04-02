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
}

func NewChangeUserRoleUseCase(userRepo ports.UserRepository) admin.ChangeUserRoleUseCase {
	return &changeUserRoleUseCase{userRepo: userRepo}
}

func (u *changeUserRoleUseCase) Execute(ctx context.Context, id string, roles []string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	rolesVo := make([]vo.Role, 0, len(roles))
	for _, r := range roles {
		validRole, err := vo.NewRole(r)
		if err != nil {
			return err
		}

		rolesVo = append(rolesVo, validRole)
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return entity.ErrUserNotFound
	}

	user.ReplaceRoles(rolesVo)

	return u.userRepo.Update(ctx, user)
}
