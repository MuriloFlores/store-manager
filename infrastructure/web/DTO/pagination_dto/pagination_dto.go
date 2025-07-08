package pagination

import "github.com/muriloFlores/StoreManager/infrastructure/web/DTO/userDTO"

type PaginationInfoResponse struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedUsersResponse struct {
	Data       []userDTO.UserResponse `json:"data"`
	Pagination PaginationInfoResponse `json:"pagination"`
}
