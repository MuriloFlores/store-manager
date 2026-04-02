package admin

import "context"

type ChangeUserRoleUseCase interface {
	Execute(ctx context.Context, id string, roles []string) error
}
