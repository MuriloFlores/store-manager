package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type FindUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewFindUserUseCase(userRepository repositories.UserRepository) *FindUserUseCase {
	return &FindUserUseCase{userRepository: userRepository}
}

func (uc *FindUserUseCase) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.userRepository.FindByEmail(ctx, email)
}

func (uc *FindUserUseCase) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return uc.userRepository.FindByID(ctx, id)
}
