package user

import "github.com/muriloFlores/StoreManager/internal/core/ports"

type UserUseCases struct {
	Create *CreateUserUseCase
	Find   *FindUserUseCase
	Update *UpdateUserUseCase
	Delete *DeleteUserUseCase
}

func NewUserUseCases(userRepo ports.UserRepository, hasher ports.PasswordHasher) *UserUseCases {
	return &UserUseCases{
		Create: nil,
		Find:   nil,
		Update: nil,
		Delete: nil,
	}
}
