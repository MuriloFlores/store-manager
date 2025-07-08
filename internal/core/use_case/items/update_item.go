package items

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type UpdateItemUseCase struct {
	itemRepo repositories.ItemRepository
	logger   ports.Logger
}

type UpdateItemParams struct {
	Name              *string
	Description       *string
	IsActive          *bool
	CanBeSold         *bool
	PriceSaleInCents  *int64
	MinimumStockLevel *int
}

func NewUpdateItemUseCase(itemRepo repositories.ItemRepository, logger ports.Logger) *UpdateItemUseCase {
	return &UpdateItemUseCase{
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (uc *UpdateItemUseCase) Execute(ctx context.Context, actor *domain.Identity, itemID string, params UpdateItemParams) (*item.Item, error) {
	uc.logger.InfoLevel("Initiate update item", map[string]interface{}{"item_ID": itemID})

	if actor.Role != value_objects.Admin && actor.Role != value_objects.Manager && actor.Role != value_objects.StockPerson {
		uc.logger.InfoLevel("user does not have permission to update item", map[string]interface{}{"id": itemID})
		return nil, &domain.ErrForbidden{Action: "trying update item"}
	}

	itemDomain, err := uc.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		uc.logger.ErrorLevel("item not found for update", err, map[string]interface{}{"id": itemID})
		return nil, err
	}

	if params.Name != nil {
		if err = itemDomain.ChangeName(*params.Name); err != nil {
			return nil, err
		}
	}

	if params.Description != nil {
		if err = itemDomain.ChangeDescription(*params.Description); err != nil {
			return nil, err
		}
	}

	if params.PriceSaleInCents != nil {
		if err = itemDomain.SetPrice(*params.PriceSaleInCents); err != nil {
			return nil, err
		}
	}

	if params.MinimumStockLevel != nil {
		if err = itemDomain.ChangeMinimumStockLevel(float64(*params.MinimumStockLevel)); err != nil {
			return nil, err
		}
	}

	if params.IsActive != nil {
		if *params.IsActive {
			itemDomain.Activate()
		} else {
			itemDomain.Deactivate()
		}
	}

	if params.CanBeSold != nil {
		itemDomain.SetCanBeSold(*params.CanBeSold)
	}

	if err = uc.itemRepo.Update(ctx, itemDomain); err != nil {
		uc.logger.ErrorLevel("failed to update item", err, map[string]interface{}{"id": itemID})
		return nil, err
	}

	uc.logger.InfoLevel("Successfully updated item", map[string]interface{}{"id": itemID})
	return itemDomain, err
}
