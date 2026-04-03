package admin

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/admin"
)

type getUsersInfoUseCase struct {
	userRepo ports.UserRepository
	logger   ports.Logger
}

func NewGetUsersInfoUseCase(userRepo ports.UserRepository, logger ports.Logger) admin.GetUsersInfo {
	return &getUsersInfoUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *getUsersInfoUseCase) Execute(ctx context.Context, pagination common.Pagination, roles []string) (*common.PaginatedResult[*entity.User], error) {
	uc.logger.Debug("fetching users info", "pagination", pagination, "filterRoles", roles)

	voRoles := make([]vo.Role, 0, len(roles))

	if len(roles) == 0 {
		voRoles = vo.AllRoles()
	} else {
		for _, role := range roles {
			validRole, err := vo.NewRole(role)
			if err != nil {
				uc.logger.Error("invalid role in filter", err, "role", role)
				return nil, err
			}
			voRoles = append(voRoles, validRole)
		}
	}

	result, err := uc.userRepo.GetUsersInfo(ctx, voRoles, pagination)
	if err != nil {
		uc.logger.Error("failed to get users info from repository", err)
		return nil, err
	}

	uc.logger.Info("users info retrieved successfully", "count", len(result.Items), "total", result.TotalCount)
	return result, nil
}
