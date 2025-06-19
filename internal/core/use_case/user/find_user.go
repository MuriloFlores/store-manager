package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type FindUserUseCase struct {
	userRepository ports.UserRepository
}

func NewFindUserUseCase(userRepository ports.UserRepository) *FindUserUseCase {
	return &FindUserUseCase{userRepository: userRepository}
}

func (uc *FindUserUseCase) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.userRepository.FindByEmail(ctx, email)
}

func (uc *FindUserUseCase) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return uc.userRepository.FindByID(ctx, id)
}
