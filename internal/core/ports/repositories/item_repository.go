package repositories

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
)

type ItemRepository interface {
	Save(ctx context.Context, item *item.Item) error
	FindByID(ctx context.Context, id string) (*item.Item, error)
	FindByIDIncludingDeleted(ctx context.Context, id string) (*item.Item, error)
	FindBySKU(ctx context.Context, sku string) (*item.Item, error)
	FindBySKUIncludingDeleted(ctx context.Context, sku string) (*item.Item, error)
	Update(ctx context.Context, item *item.Item) error
	Delete(ctx context.Context, itemID string) error
	List(ctx context.Context, params *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error)
	ListForUsers(ctx context.Context, paginationParams *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error)
	Search(ctx context.Context, searchTerm string, isPublicSearch bool, params *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error)
}
