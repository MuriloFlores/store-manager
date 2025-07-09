package items

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type FindItemUseCase struct {
	itemRepo repositories.ItemRepository
	logger   ports.Logger
}

func NewFindItemUseCase(
	itemRepo repositories.ItemRepository,
	logger ports.Logger,
) *FindItemUseCase {
	return &FindItemUseCase{
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (uc *FindItemUseCase) FindByID(ctx context.Context, id string, actor *domain.Identity) (*item.Item, error) {
	uc.logger.InfoLevel("Invoke Find By ID", map[string]interface{}{"id": id})

	if actor.Role.IsStockEmployee() {
		uc.logger.InfoLevel("user not allowed")
		return nil, &domain.ErrForbidden{Action: "You don't have permission to create an item"}
	}

	return uc.itemRepo.FindByID(ctx, id)
}

func (uc *FindItemUseCase) FindBySKU(ctx context.Context, sku string, actor *domain.Identity) (*item.Item, error) {
	uc.logger.InfoLevel("Invoke Find By SKU", map[string]interface{}{"sku": sku})

	if actor.Role.IsStockEmployee() {
		uc.logger.InfoLevel("user not allowed")
		return nil, &domain.ErrForbidden{Action: "You don't have permission to create an item"}
	}

	return uc.itemRepo.FindBySKU(ctx, sku)
}
