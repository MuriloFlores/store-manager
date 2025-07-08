package items

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type DeleteItemUseCase struct {
	itemRepo repositories.ItemRepository
	userRepo repositories.UserRepository
	logger   ports.Logger
}

func NewDeleteItemUseCase(itemRepo repositories.ItemRepository, logger ports.Logger) *DeleteItemUseCase {
	return &DeleteItemUseCase{
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (uc *DeleteItemUseCase) Execute(ctx context.Context, actor *domain.Identity, itemID string) error {
	uc.logger.InfoLevel("Init delete item", map[string]interface{}{"item_ID": itemID})

	if actor.Role != value_objects.Admin && actor.Role != value_objects.Manager && actor.Role != value_objects.StockPerson {
		return &domain.ErrForbidden{Action: "user does not have permission to delete item"}
	}

	_, err := uc.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		uc.logger.ErrorLevel("failed in find item", err, map[string]interface{}{"item_ID": itemID})
		return err
	}

	if err = uc.itemRepo.Delete(ctx, itemID); err != nil {
		uc.logger.ErrorLevel("failed in delete item", err, map[string]interface{}{"item_ID": itemID})
		return err
	}

	return nil
}
