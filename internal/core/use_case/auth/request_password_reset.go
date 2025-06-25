package auth

import (
	"context"
	"errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"time"
)

type RequestPasswordResetUseCase struct {
	userRepo       repositories.UserRepository
	tokenRepo      repositories.ActionTokenRepository
	tokenGenerator ports.SecureTokenGenerator
	taskEnqueuer   ports.TaskEnqueuer
}

func NewRequestPasswordResetUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.ActionTokenRepository,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		tokenGenerator: tokenGenerator,
		taskEnqueuer:   taskEnqueuer,
	}
}

func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, email string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		var notFoundErr *domain.ErrNotFound
		if errors.As(err, &notFoundErr) {
			//adicionar logger posteriormente
			return nil
		}

		return err
	}

	resetTokenString, err := uc.tokenGenerator.Generate()
	if err != nil {
		return err
	}

	resetToken := &domain.ActionToken{
		Token:     resetTokenString,
		UserID:    user.ID(),
		Type:      domain.PasswordReset,
		Payload:   "",
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	if err = uc.tokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	jobData := jobs.PasswordResetJobData{
		UserName:  user.Name(),
		UserEmail: user.Email(),
		ResetLink: "https://app.muriloflores.xyz/reset-password?token=" + resetToken.Token,
	}

	if err = uc.taskEnqueuer.EnqueuePasswordReset(&jobData); err != nil {
		return err
	}

	return nil
}
