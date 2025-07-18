package repositories

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID string) error
	CountAdmins(ctx context.Context) (int, error)
	FindByEmailIncludingDeleted(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, params *pagination.PaginationParams) (*pagination.PaginatedResult[*domain.User], error)
	Search(ctx context.Context, searchTerm string, params *pagination.PaginationParams) (*pagination.PaginatedResult[*domain.User], error)
}
