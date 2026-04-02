package admin

import "context"

type ChangeUserStatusUseCase interface {
	Execute(ctx context.Context, id string, active bool) error
}
