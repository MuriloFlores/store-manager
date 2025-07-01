package items

import (
	"context"
	"errors"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type CreateItemParams struct {
	Name              string
	Description       string
	SKU               string
	ItemType          item.ItemType
	IsActive          bool
	CanBeSold         bool
	PriceSaleInCents  int64
	StockQuantity     float64
	UnitOfMeasure     string
	MinimumStockLevel float64
}

type CreateItemUseCase struct {
	itemRepo  repositories.ItemRepository
	logger    ports.Logger
	generator ports.IDGenerator
}

func NewCreateItemUseCase(
	itemRepo repositories.ItemRepository,
	logger ports.Logger,
	generator ports.IDGenerator,
) *CreateItemUseCase {
	return &CreateItemUseCase{
		itemRepo:  itemRepo,
		logger:    logger,
		generator: generator,
	}
}

func (uc *CreateItemUseCase) Execute(ctx context.Context, params CreateItemParams) (*item.Item, error) {
	uc.logger.InfoLevel("Invoke Create Item Use Case", map[string]interface{}{"item": ""})

	if params.SKU != "" {
		existing, err := uc.itemRepo.FindBySKU(ctx, params.SKU)
		if err != nil {
			var notFoundErr *domain.ErrNotFound
			if !errors.As(err, &notFoundErr) {
				uc.logger.ErrorLevel("failed to check for existing SKU", err, map[string]interface{}{"item": params.Name, "sku": params.SKU})
				return nil, err
			}
		}

		if existing != nil {
			return nil, &domain.ErrConflict{Resource: "item", Details: fmt.Sprintf("SKU %s already exists", params.SKU)}
		}
	}

	id := uc.generator.Generate()

	itemDomain, err := item.NewItemBuilder().
		WithID(id).
		WithName(params.Name).
		WithSKU(params.SKU).
		WithDescription(params.Description).
		WithType(params.ItemType).
		WithPriceInCents(params.PriceSaleInCents).
		WithUnitOfMeasure(params.UnitOfMeasure).
		WithQuantity(params.StockQuantity).
		WithMinimumStock(params.MinimumStockLevel).
		WithCanBeSold(params.CanBeSold).
		Build()

	if err != nil {
		uc.logger.ErrorLevel("failed to create item domain", err, map[string]interface{}{"item": params.Name, "sku": params.SKU})
		return nil, err
	}

	err = uc.itemRepo.Save(ctx, itemDomain)
	if err != nil {
		return nil, err
	}

	uc.logger.InfoLevel("item saved successfully", map[string]interface{}{"item": params.Name, "sku": params.SKU})

	return itemDomain, nil
}
