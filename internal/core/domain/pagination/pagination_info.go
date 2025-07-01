package pagination

import (
	"math"
)

type PaginationInfo struct {
	CurrentPage int
	PageSize    int
	TotalItems  int64
	TotalPages  int
}

func (p *PaginationInfo) CalculateTotalPages() {
	if p.PageSize == 0 {
		p.TotalPages = 0
		return
	}
	p.TotalPages = int(math.Ceil(float64(p.TotalItems) / float64(p.PageSize)))
}
