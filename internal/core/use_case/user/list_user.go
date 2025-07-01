package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type ListUsersUseCase struct {
	userRepo repositories.UserRepository
	logger   ports.Logger
}

func NewListUsersUseCase(userRepo repositories.UserRepository, logger ports.Logger) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, actor *domain.Identity, params *pagination.PaginationParams) (*pagination.PaginatedResult[*domain.User], error) {
	uc.logger.InfoLevel("Executing ListUsersUseCase", map[string]interface{}{"actor_id": actor.UserID, "params": params})

	if actor.Role != value_objects.Admin && actor.Role != value_objects.Manager {
		uc.logger.InfoLevel("Unauthorized access attempt", map[string]interface{}{"actor_id": actor.UserID, "role": actor.Role})
		return nil, &domain.ErrForbidden{Action: "list users"}
	}

	return uc.userRepo.List(ctx, params)
}
