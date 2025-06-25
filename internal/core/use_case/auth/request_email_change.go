package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"time"
)

type RequestEmailChangeUseCase struct {
	userRepo       repositories.UserRepository
	tokenRepo      repositories.ActionTokenRepository
	tokenGenerator ports.SecureTokenGenerator
	hasher         ports.PasswordHasher
	taskEnqueuer   ports.TaskEnqueuer
}

func NewRequestEmailChangeUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.ActionTokenRepository,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	hasher ports.PasswordHasher,
) *RequestEmailChangeUseCase {
	return &RequestEmailChangeUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		tokenGenerator: tokenGenerator,
		taskEnqueuer:   taskEnqueuer,
		hasher:         hasher,
	}
}

func (uc *RequestEmailChangeUseCase) Execute(ctx context.Context, actor *domain.Identity, currentPassword, newEmail string) error {
	user, err := uc.userRepo.FindByID(ctx, actor.UserID)
	if err != nil {
		return err
	}

	if !uc.hasher.Compare(user.Password(), currentPassword) {
		return &domain.ErrInvalidCredentials{}
	}

	if newEmail == "" {
		return &domain.ErrInvalidInput{FieldName: "newEmail", Reason: "new email can't be empty"}
	}

	if newEmail == user.Email() {
		return &domain.ErrForbidden{Action: "the emails must be different"}
	}

	existingUser, err := uc.userRepo.FindByID(ctx, newEmail)
	if err == nil && existingUser != nil {
		return &domain.ErrConflict{Resource: "email", Details: "the email provided is already in use"}
	}

	var notFoundErr *domain.ErrNotFound
	if err != nil && !errors.As(err, &notFoundErr) {
		return err
	}

	confirmationTokenString, err := uc.tokenGenerator.Generate()
	if err != nil {
		return err
	}

	changeToken := &domain.ActionToken{
		Token:     confirmationTokenString,
		UserID:    user.ID(),
		Type:      domain.EmailConfirmation,
		Payload:   newEmail,
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	if err = uc.tokenRepo.Create(ctx, changeToken); err != nil {
		return err
	}

	confirmationJob := jobs.EmailChangeConfirmationJobData{
		UserName:         user.Name(),
		ConfirmationLink: "https://app.muriloflores.xyz/auth/confirm-email?token=" + changeToken.Token,
		ToEmail:          user.Email(),
	}

	if err = uc.taskEnqueuer.EnqueueEmailChangeConfirmation(&confirmationJob); err != nil {
		return fmt.Errorf("internal error starting email change process: %w", err)
	}

	alertJob := &jobs.SecurityNotificationJobData{
		UserName: user.Name(),
		ToEmail:  user.Email(),
		Message: fmt.Sprintf(`Recebemos uma solicitação para alterar seu email para 
		'%s'. Se você não fez essa solicitação, por favor, ignore este email.`, newEmail),
	}

	if err = uc.taskEnqueuer.EnqueueSecurityNotification(alertJob); err != nil {
		return fmt.Errorf("internal error when starting the confirmation email sending process: %w", err)
	}

	return nil
}
