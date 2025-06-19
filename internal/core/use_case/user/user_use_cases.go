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

func NewUserUseCases(userRepo ports.UserRepository, hasher ports.PasswordHasher, generator ports.IDGenerator) *UserUseCases {
	return &UserUseCases{
		Create: NewCreateUserUseCase(userRepo, hasher, generator),
		Find:   NewFindUserUseCase(userRepo),
		Update: NewUpdateUserUseCase(userRepo, hasher),
		Delete: NewDeleteUserUseCase(userRepo),
	}
}
