package user

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/google/uuid"
)

type MyInfoUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error)
}
