package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"strings"
)

type SearchUsersUseCase struct {
	userRepo repositories.UserRepository
	logger   ports.Logger
}

func NewSearchUsersUseCase(
	userRepo repositories.UserRepository,
	logger ports.Logger,
) *SearchUsersUseCase {
	return &SearchUsersUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *SearchUsersUseCase) Execute(ctx context.Context, actor *domain.Identity, searchTerm string, params *pagination.PaginationParams) (*pagination.PaginatedResult[*domain.User], error) {
	uc.logger.InfoLevel("Invoking Search Users Use Case")

	if actor == nil || !actor.Role.IsAdminEmployee() {
		uc.logger.InfoLevel("user not allowed to search users")
		return nil, &domain.ErrForbidden{}
	}

	return uc.userRepo.Search(ctx, strings.TrimSpace(searchTerm), params)
}
