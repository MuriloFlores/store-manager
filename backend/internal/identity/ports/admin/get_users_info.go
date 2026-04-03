package admin

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
)

type GetUsersInfo interface {
	Execute(ctx context.Context, pagination common.Pagination, roles []string) (*common.PaginatedResult[*entity.User], error)
}
