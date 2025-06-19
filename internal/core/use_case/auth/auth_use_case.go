package auth

import "github.com/muriloFlores/StoreManager/internal/core/ports"

type AuthUseCases struct {
	Login                *LoginUserUseCase
	ChangePassword       *ChangePasswordUseCase
	RequestPasswordReset *RequestPasswordResetUseCase
	ConfirmPasswordReset *ConfirmPasswordResetUseCase
	RequestEmailChange   *RequestEmailChangeUseCase
	ConfirmEmailChange   *ConfirmEmailChangeUseCase
}

func NewAuthUseCases(
	userRepo ports.UserRepository,
	hasher ports.PasswordHasher,
	manager ports.TokenManager,
	tokenRepo ports.ActionTokenRepository,
	tokenGen ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
) *AuthUseCases {
	return &AuthUseCases{
		Login:                NewLoginUserUseCase(userRepo, hasher, manager),
		ChangePassword:       NewChangePasswordUseCase(userRepo, hasher),
		RequestPasswordReset: NewRequestPasswordResetUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer),
		ConfirmPasswordReset: NewConfirmPasswordResetUseCase(userRepo, tokenRepo, hasher),
		RequestEmailChange:   NewRequestEmailChangeUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer, hasher),
		ConfirmEmailChange:   NewConfirmEmailChangeUseCase(userRepo, tokenRepo),
	}
}
