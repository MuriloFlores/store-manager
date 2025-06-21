package user

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
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
	logger ports.Logger,
	accountValidation auth.RequestAccountValidationUseCase,

) *UserUseCases {
	return &UserUseCases{
		Create: NewCreateUserUseCase(userRepo, hasher, generator, tokenGenerator, taskEnqueuer, tokenRepo, logger, accountValidation),
		Find:   NewFindUserUseCase(userRepo),
		Update: NewUpdateUserUseCase(userRepo, hasher),
		Delete: NewDeleteUserUseCase(userRepo),
	}
}
