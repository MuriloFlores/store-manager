package admin

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
)

type GetUsersInfo interface {
	Execute(ctx context.Context, pagination _common.Pagination, roles []string) (*_common.PaginatedResult[*entity.User], error)
}
