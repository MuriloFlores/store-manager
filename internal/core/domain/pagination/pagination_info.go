package domain

import (
	"math"
)

type PaginationParams struct {
	Page     int
	PageSize int
}

type PaginationInfo struct {
	CurrentPage int
	PageSize    int
	TotalItems  int64
	TotalPages  int
}

type PaginatedEntity interface {
	*User | *Item
}

type PaginatedResult[T PaginatedEntity] struct {
	Data       []T
	Pagination PaginationInfo
}

func (p *PaginationInfo) CalculateTotalPages() {
	if p.PageSize == 0 {
		p.TotalPages = 0
		return
	}
	p.TotalPages = int(math.Ceil(float64(p.TotalItems) / float64(p.PageSize)))
}
