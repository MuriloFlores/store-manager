package auth

import "github.com/muriloFlores/StoreManager/internal/core/ports"

type AuthUseCases struct {
	Login                     *LoginUserUseCase
	ChangePassword            *ChangePasswordUseCase
	ConfirmPasswordReset      *ConfirmPasswordResetUseCase
	ConfirmUserEmailUseCase   *ConfirmUserEmailUseCase
	ConfirmEmailChangeUseCase *ConfirmEmailChangeUseCase
	RequestEmailChange        *RequestEmailChangeUseCase
	RequestPasswordReset      *RequestPasswordResetUseCase
	ConfirmAccountUserUseCase *ConfirmAccountUserUseCase
}

func NewAuthUseCases(
	userRepo ports.UserRepository,
	hasher ports.PasswordHasher,
	manager ports.TokenManager,
	tokenRepo ports.ActionTokenRepository,
	tokenGen ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	logger ports.Logger,
) *AuthUseCases {
	return &AuthUseCases{
		Login:                     NewLoginUserUseCase(userRepo, hasher, manager, logger),
		ChangePassword:            NewChangePasswordUseCase(userRepo, hasher),
		RequestPasswordReset:      NewRequestPasswordResetUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer),
		ConfirmPasswordReset:      NewConfirmPasswordResetUseCase(userRepo, tokenRepo, hasher),
		RequestEmailChange:        NewRequestEmailChangeUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer, hasher),
		ConfirmEmailChangeUseCase: NewConfirmEmailChangeUseCase(userRepo, tokenRepo),
		ConfirmUserEmailUseCase:   NewConfirmUserEmailUseCase(userRepo, tokenRepo, logger),
		ConfirmAccountUserUseCase: NewConfirmAccountUserUseCase(userRepo, tokenRepo, logger),
	}
}
