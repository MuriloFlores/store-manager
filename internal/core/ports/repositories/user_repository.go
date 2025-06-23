package repositories

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID string) error
	CountAdmins(ctx context.Context) (int, error)
	FindByEmailIncludingDeleted(ctx context.Context, email string) (*domain.User, error)
}
