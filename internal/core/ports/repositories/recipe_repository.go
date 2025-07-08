package repositories

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
)

type RecipeRepository interface {
	Save(ctx context.Context, recipe *domain.Recipe) error
	FindByProductID(ctx context.Context, productID string) (*domain.Recipe, error)
}
