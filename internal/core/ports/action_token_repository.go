package ports

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
)

type ActionTokenRepository interface {
	Create(ctx context.Context, token *domain.ActionToken) error
	FindAndConsume(ctx context.Context, tokenString string, tokenType domain.ActionType) (*domain.ActionToken, error)
}
