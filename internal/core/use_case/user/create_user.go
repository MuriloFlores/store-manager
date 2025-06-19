package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
	"time"
)

type CreateUserUseCase struct {
	userRepository ports.UserRepository
	hasher         ports.PasswordHasher
	generator      ports.IDGenerator
	tokenGenerator ports.SecureTokenGenerator
	taskEnqueuer   ports.TaskEnqueuer
	tokenRepo      ports.ActionTokenRepository
}

func NewCreateUserUseCase(userRepository ports.UserRepository, hasher ports.PasswordHasher, generator ports.IDGenerator) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
		hasher:         hasher,
		generator:      generator,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, name, email, password string, role value_objects.Role) (*domain.User, error) {
	_, err := uc.userRepository.FindByEmail(ctx, email)
	if err == nil {
		return nil, &domain.ErrConflict{Resource: "user", Details: "email already used in another user "}
	}

	var notFoundErr *domain.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		return nil, fmt.Errorf("error verifying user email: %w", err)
	}

	hashedPassword, err := uc.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	id := uc.generator.Generate()

	user, err := domain.NewUser(id, name, email, hashedPassword, role)
	if err != nil {
		return nil, err
	}

	if err = uc.userRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	verificationTokenString, err := uc.tokenGenerator.Generate()
	if err != nil {
		return nil, err
	}

	actionToken := &domain.ActionToken{
		Token:     verificationTokenString,
		UserID:    user.ID(),
		Type:      domain.EmailConfirmation,
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	if err = uc.tokenRepo.Create(ctx, actionToken); err != nil {
		return nil, err
	}

	jobData := &domain.EmailChangeConfirmationJobData{
		UserName:         user.Name(),
		ToEmail:          user.Email(),
		ConfirmationLink: "http://localhost/verify-account?token=" + actionToken.Token,
	}

	if err = uc.taskEnqueuer.EnqueueEmailChangeConfirmation(jobData); err != nil {
		return nil, err
	}

	return user, nil
}
