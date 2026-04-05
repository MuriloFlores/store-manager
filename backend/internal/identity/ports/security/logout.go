package security

import "context"

type LogoutUseCase interface {
	Execute(ctx context.Context, refreshToken string) error
}
