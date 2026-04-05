package security

import (
	"context"

	"github.com/google/uuid"
)

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
}
