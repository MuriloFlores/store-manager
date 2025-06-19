package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"time"
)

type RequestEmailChangeUseCase struct {
	userRepo       ports.UserRepository
	tokenRepo      ports.ActionTokenRepository
	tokenGenerator ports.SecureTokenGenerator
	hasher         ports.PasswordHasher
	taskEnqueuer   ports.TaskEnqueuer
}

func NewRequestEmailChangeUseCase(
	userRepo ports.UserRepository,
	tokenRepo ports.ActionTokenRepository,
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

	confirmationJob := domain.EmailChangeConfirmationJobData{
		UserName:         user.Name(),
		ConfirmationLink: "http://localhost/change-email?token=" + changeToken.Token,
		ToEmail:          user.Email(),
	}

	if err = uc.taskEnqueuer.EnqueueEmailChangeConfirmation(ctx, &confirmationJob); err != nil {
		// deveriamos remover o token se a fila falhar?
	}

	alertJob := &domain.SecurityNotificationJobData{
		UserName: user.Name(),
		ToEmail:  user.Email(),
		Message: fmt.Sprintf(`Recebemos uma solicitação para alterar seu email para 
		'%s'. Se você não fez essa solicitação, por favor, ignore este email.`, newEmail),
	}

	if err = uc.taskEnqueuer.EnqueueSecurityNotification(ctx, alertJob); err != nil {
		// Se a fila falhar, deveriamos remover o token?
	}

	//pq nao usamos mais o enqueue?
	return nil
}
