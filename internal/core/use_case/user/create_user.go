package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type CreateUserUseCase struct {
	userRepository    ports.UserRepository
	hasher            ports.PasswordHasher
	generator         ports.IDGenerator
	tokenGenerator    ports.SecureTokenGenerator
	taskEnqueuer      ports.TaskEnqueuer
	tokenRepo         ports.ActionTokenRepository
	logger            ports.Logger
	accountValidation auth.RequestAccountValidationUseCase //não tenho certeza se isso foi importado corretamente, no sentido de uma arquitetura limpa, talvez seja melhor criar uma nova port para isso
}

func NewCreateUserUseCase(
	userRepository ports.UserRepository,
	hasher ports.PasswordHasher,
	generator ports.IDGenerator,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	tokenRepo ports.ActionTokenRepository,
	logger ports.Logger,
	accountValidation auth.RequestAccountValidationUseCase,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository:    userRepository,
		hasher:            hasher,
		generator:         generator,
		tokenGenerator:    tokenGenerator,
		taskEnqueuer:      taskEnqueuer,
		tokenRepo:         tokenRepo,
		logger:            logger,
		accountValidation: accountValidation,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, name, email, password string, role value_objects.Role) (*domain.User, error) {
	uc.logger.InfoLevel("Create user use case started")

	_, err := uc.userRepository.FindByEmail(ctx, email)
	if err == nil {
		uc.logger.ErrorLevel("Email already used in another user", err, map[string]interface{}{"user_email": email})
		return nil, &domain.ErrConflict{Resource: "user", Details: "email already used in another user "}
	}

	var notFoundErr *domain.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		uc.logger.ErrorLevel("Error finding user by email", err, map[string]interface{}{"user_email": email})
		return nil, fmt.Errorf("error verifying user email: %w", err)
	}

	hashedPassword, err := uc.hasher.Hash(password)
	if err != nil {
		uc.logger.ErrorLevel("Error hashing password", err, map[string]interface{}{"user_name": name, "user_email": email})
		return nil, err
	}

	id := uc.generator.Generate()

	user, err := domain.NewUser(id, name, email, hashedPassword, role)
	if err != nil {
		uc.logger.ErrorLevel("Error creating new user", err, map[string]interface{}{"user_id": id})
		return nil, err
	}

	if err = uc.userRepository.Save(ctx, user); err != nil {
		uc.logger.ErrorLevel("Error saving user to repository", err, map[string]interface{}{"user_id": id})
		return nil, err
	}

	if err = uc.accountValidation.Execute(ctx, user.ID()); err != nil {
		uc.logger.ErrorLevel("Error requesting account validation", err, map[string]interface{}{"user_id": id})
		return nil, err
	}

	uc.logger.InfoLevel("User created successfully", map[string]interface{}{"user_id": id, "user_name": user.Name()})
	return user, nil
}
