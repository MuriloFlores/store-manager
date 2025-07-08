package pagination_dto

import (
	"github.com/muriloFlores/StoreManager/infrastructure/web/DTO/item_dto"
	"github.com/muriloFlores/StoreManager/infrastructure/web/DTO/user_dto"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
	"net/http"
	"strconv"
)

type PaginationInfoResponse struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedDTO interface {
	user_dto.UserResponse | item_dto.InternalItemResponse | item_dto.ClientItemResponse
}

type PaginatedResponse[T PaginatedDTO] struct {
	Data       []T                    `json:"data"`
	Pagination PaginationInfoResponse `json:"pagination"`
}

func ParsePagination(r *http.Request) *pagination.PaginationParams {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("page_size")

	page, _ := strconv.Atoi(pageQuery)
	pageSize, _ := strconv.Atoi(pageSizeQuery)

	if page < 1 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	return &pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

func ToPaginationInfoResponse(domainPagination pagination.PaginationInfo) PaginationInfoResponse {
	return PaginationInfoResponse{
		CurrentPage: domainPagination.CurrentPage,
		PageSize:    domainPagination.PageSize,
		TotalItems:  domainPagination.TotalItems,
		TotalPages:  domainPagination.TotalPages,
	}
}
