package admin

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/admin"
)

type getUsersInfoUseCase struct {
	userRepo ports.UserRepository
}

func NewGetUsersInfoUseCase(userRepo ports.UserRepository) admin.GetUsersInfo {
	return &getUsersInfoUseCase{userRepo: userRepo}
}

func (uc *getUsersInfoUseCase) Execute(ctx context.Context, pagination _common.Pagination, roles []string) (*_common.PaginatedResult[*entity.User], error) {
	voRoles := make([]vo.Role, 0, len(roles))

	if len(roles) == 0 {
		voRoles = vo.AllRoles()
	}

	for _, role := range roles {
		validRole, err := vo.NewRole(role)
		if err != nil {
			return nil, err
		}
		voRoles = append(voRoles, validRole)
	}

	return uc.userRepo.GetUsersInfo(ctx, voRoles, pagination)
}
