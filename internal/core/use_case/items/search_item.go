package items

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type SearchItemUseCase struct {
	itemRepo repositories.ItemRepository
	logger   ports.Logger
}

func NewSearchItemUseCase(
	itemRepo repositories.ItemRepository,
	logger ports.Logger,
) *SearchItemUseCase {
	return &SearchItemUseCase{
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (uc *SearchItemUseCase) Execute(ctx context.Context, actor *domain.Identity, searchTerm string, params *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	uc.logger.InfoLevel("Invoking Search Use Case")

	isPublicSearch := true

	if actor.Role.IsStockEmployee() {
		uc.logger.InfoLevel("user not allowed")
		isPublicSearch = false
	}

	return uc.itemRepo.Search(ctx, searchTerm, isPublicSearch, params)
}
