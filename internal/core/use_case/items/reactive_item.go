package items

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type ReactiveItemUseCase struct {
	itemRepo repositories.ItemRepository
	logger   ports.Logger
}

func NewReactiveItemUseCase(itemRepo repositories.ItemRepository, logger ports.Logger) *ReactiveItemUseCase {
	return &ReactiveItemUseCase{
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (uc *ReactiveItemUseCase) Execute(ctx context.Context, actor *domain.Identity, itemID string) (*item.Item, error) {
	uc.logger.InfoLevel("Starting item reactivation process")

	if !actor.Role.IsStockEmployee() {
		return nil, &domain.ErrForbidden{Action: "reactivate Item"}
	}

	itemDomain, err := uc.itemRepo.FindByIDIncludingDeleted(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if !itemDomain.IsDeleted() {
		return nil, &domain.ErrInvalidInput{Reason: "item is not deleted"}
	}

	itemDomain.SetDeleted(nil)

	if err = uc.itemRepo.Update(ctx, itemDomain); err != nil {
		uc.logger.ErrorLevel("failed to update item for reactivation", err)
		return nil, err
	}

	return itemDomain, nil
}
