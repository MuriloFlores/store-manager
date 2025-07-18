package auth

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
)

type AuthUseCases struct {
	Login                           *LoginUserUseCase
	ChangePassword                  *ChangePasswordUseCase
	ConfirmPasswordReset            *ConfirmPasswordResetUseCase
	ConfirmUserEmailUseCase         *ConfirmUserEmailUseCase
	ConfirmEmailChangeUseCase       *ConfirmEmailChangeUseCase
	RequestEmailChange              *RequestEmailChangeUseCase
	RequestPasswordReset            *RequestPasswordResetUseCase
	ConfirmAccountUserUseCase       *ConfirmAccountUserUseCase
	RequestAccountValidationUseCase *RequestAccountValidationUseCase
}

func NewAuthUseCases(
	userRepo repositories.UserRepository,
	hasher ports.PasswordHasher,
	manager ports.TokenManager,
	tokenRepo repositories.ActionTokenRepository,
	tokenGen ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	logger ports.Logger,
	limiter ports.RateLimiter,
) *AuthUseCases {
	return &AuthUseCases{
		Login:                           NewLoginUserUseCase(userRepo, hasher, manager, logger),
		ChangePassword:                  NewChangePasswordUseCase(userRepo, hasher),
		RequestPasswordReset:            NewRequestPasswordResetUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer),
		ConfirmPasswordReset:            NewConfirmPasswordResetUseCase(userRepo, tokenRepo, hasher),
		RequestEmailChange:              NewRequestEmailChangeUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer, hasher),
		ConfirmEmailChangeUseCase:       NewConfirmEmailChangeUseCase(userRepo, tokenRepo),
		ConfirmUserEmailUseCase:         NewConfirmUserEmailUseCase(userRepo, tokenRepo, logger),
		ConfirmAccountUserUseCase:       NewConfirmAccountUserUseCase(userRepo, tokenRepo, logger),
		RequestAccountValidationUseCase: NewRequestAccountValidationUseCase(userRepo, tokenRepo, tokenGen, taskEnqueuer, logger, limiter),
	}
}
