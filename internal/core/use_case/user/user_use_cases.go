package user

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type UserUseCases struct {
	Create *CreateUserUseCase
	Find   *FindUserUseCase
	Update *UpdateUserUseCase
	Delete *DeleteUserUseCase
}

func NewUserUseCases(
	userRepo ports.UserRepository,
	hasher ports.PasswordHasher,
	generator ports.IDGenerator,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	tokenRepo ports.ActionTokenRepository,
) *UserUseCases {
	return &UserUseCases{
		Create: NewCreateUserUseCase(userRepo, hasher, generator, tokenGenerator, taskEnqueuer, tokenRepo),
		Find:   NewFindUserUseCase(userRepo),
		Update: NewUpdateUserUseCase(userRepo, hasher),
		Delete: NewDeleteUserUseCase(userRepo),
	}
}
