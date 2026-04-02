package auth

import (
	"context"
)

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, email string) error
}
