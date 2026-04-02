package auth

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
)

type LoginUseCase interface {
	Execute(ctx context.Context, input *dto.LoginRequest) (*dto.LoginResult, error)
}
