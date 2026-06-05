package ports

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/organization/domain/entity"
)

type StoreRepository interface {
	Save(ctx context.Context, store *entity.Store) error
}
