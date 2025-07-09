package items

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type ListItemUseCase struct {
	itemRepo repositories.ItemRepository
	logger   ports.Logger
}

func NewListItemUseCase(itemRepo repositories.ItemRepository, logger ports.Logger) *ListItemUseCase {
	return &ListItemUseCase{itemRepo: itemRepo, logger: logger}
}

func (uc *ListItemUseCase) ListPublic(ctx context.Context, params *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	uc.logger.InfoLevel("Executing ListUserUseCase")

	return uc.itemRepo.ListForUsers(ctx, params)
}

func (uc *ListItemUseCase) ListInternal(ctx context.Context, actor *domain.Identity, params *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	uc.logger.InfoLevel("executing list internal items use case", map[string]interface{}{"actor_id": actor.UserID})

	if !actor.Role.IsStockEmployee() {
		return nil, &domain.ErrForbidden{Action: "list internal items"}
	}

	uc.logger.InfoLevel("User has permission, fetching internal item list", map[string]interface{}{"actor_id": actor.UserID})
	return uc.itemRepo.List(ctx, params)
}
