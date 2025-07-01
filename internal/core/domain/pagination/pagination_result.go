package pagination

import (
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
)

type PaginatedEntity interface {
	*domain.User | *item.Item
}

type PaginatedResult[T PaginatedEntity] struct {
	Data       []T
	Pagination PaginationInfo
}
