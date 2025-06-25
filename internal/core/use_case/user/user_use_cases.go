package user

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
)

type UserUseCases struct {
	Create  *CreateClienteUseCase
	Find    *FindUserUseCase
	Update  *UpdateUserUseCase
	Delete  *DeleteUserUseCase
	Promote *PromoteUserUseCase
	List    *ListUsersUseCase
}

func NewUserUseCases(
	userRepo repositories.UserRepository,
	hasher ports.PasswordHasher,
	generator ports.IDGenerator,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	tokenRepo repositories.ActionTokenRepository,
	logger ports.Logger,
	accountValidation auth.RequestAccountValidationUseCase,

) *UserUseCases {
	return &UserUseCases{
		Create:  NewCreateClienteUseCase(userRepo, hasher, generator, tokenGenerator, taskEnqueuer, tokenRepo, logger, accountValidation),
		Find:    NewFindUserUseCase(userRepo),
		Update:  NewUpdateUserUseCase(userRepo, hasher),
		Delete:  NewDeleteUserUseCase(userRepo),
		Promote: NewPromoteUseCase(userRepo, logger, taskEnqueuer),
		List:    NewListUsersUseCase(userRepo, logger),
	}
}
