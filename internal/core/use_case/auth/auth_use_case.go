package auth

import "github.com/muriloFlores/StoreManager/internal/core/ports"

type AuthUseCases struct {
	Login          *LoginUserUseCase
	ChangePassword *ChangePasswordUseCase
	//RequestPasswordReset
	//ConfirmPasswordReset
	//RequestEmailChange
	//ConfirmEmailChange
}

func NewAuthUseCases(repository ports.UserRepository, hasher ports.PasswordHasher, manager ports.TokenManager) *AuthUseCases {
	return &AuthUseCases{
		Login:          NewLoginUserUseCase(repository, hasher, manager),
		ChangePassword: NewChangePasswordUseCase(repository, hasher),
	}
}
