package items

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type ItemsUseCases struct {
	Create *CreateItemUseCase
	Find   *FindItemUseCase
	Update *UpdateItemUseCase
	Delete *DeleteItemUseCase
	List   *ListItemUseCase
}

func NewItemUseCases(
	itemRepo repositories.ItemRepository,
	logger ports.Logger,
	generator ports.IDGenerator,
) *ItemsUseCases {
	return &ItemsUseCases{
		Create: NewCreateItemUseCase(itemRepo, logger, generator),
		Find:   NewFindItemUseCase(itemRepo, logger),
		Update: NewUpdateItemUseCase(itemRepo, logger),
		Delete: NewDeleteItemUseCase(itemRepo, logger),
		List:   NewListItemUseCase(itemRepo, logger),
	}
}
