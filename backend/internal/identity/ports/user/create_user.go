package user

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, input dto.CreateUserInput) error
}
