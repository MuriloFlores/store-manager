package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type CreateClienteUseCase struct {
	userRepository    repositories.UserRepository
	hasher            ports.PasswordHasher
	generator         ports.IDGenerator
	tokenGenerator    ports.SecureTokenGenerator
	taskEnqueuer      ports.TaskEnqueuer
	tokenRepo         repositories.ActionTokenRepository
	logger            ports.Logger
	accountValidation auth.RequestAccountValidationUseCase //não tenho certeza se isso foi importado corretamente, no sentido de uma arquitetura limpa, talvez seja melhor criar uma nova port para isso
}

func NewCreateClienteUseCase(
	userRepository repositories.UserRepository,
	hasher ports.PasswordHasher,
	generator ports.IDGenerator,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	tokenRepo repositories.ActionTokenRepository,
	logger ports.Logger,
	accountValidation auth.RequestAccountValidationUseCase,
) *CreateClienteUseCase {
	return &CreateClienteUseCase{
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

func (uc *CreateClienteUseCase) Execute(ctx context.Context, name, email, password string) (*domain.User, error) {
	uc.logger.InfoLevel("Creating client", map[string]interface{}{"name": name, "email": email})

	user, err := uc.userRepository.FindByEmailIncludingDeleted(ctx, email)
	if err == nil {
		if !user.IsDeleted() {
			uc.logger.ErrorLevel("User already exists", errors.New("user already exists"), map[string]interface{}{"email": email})
			return nil, &domain.ErrConflict{Resource: "email", Details: "User already exists with this email"}
		}

		if err = uc.reactivateUser(ctx, user, name, password, value_objects.Client); err != nil {
			uc.logger.ErrorLevel("Error reactivating user", err, map[string]interface{}{"user_id": user.ID()})
			return nil, err
		}

		uc.logger.InfoLevel("User reactivated successfully", map[string]interface{}{"user_id": user.ID()})

	} else {
		var notFoundErr *domain.ErrNotFound
		if !errors.As(err, &notFoundErr) {
			uc.logger.ErrorLevel("Error finding user by email", err, map[string]interface{}{"email": email})
			return nil, &domain.ErrInvalidInput{}
		}

		userID := uc.generator.Generate()
		hashedPassword, err := uc.hasher.Hash(password)
		if err != nil {
			uc.logger.ErrorLevel("Error hashing password", err, map[string]interface{}{"email": email})
			return nil, &domain.ErrInvalidInput{FieldName: "password", Reason: "failed to hash password"}
		}

		domainUser, err := domain.NewUser(userID, name, email, hashedPassword, value_objects.Client)
		if err != nil {
			uc.logger.ErrorLevel("Error creating new user", err, map[string]interface{}{"email": email})
			return nil, err
		}

		err = uc.userRepository.Save(ctx, domainUser)
		if err != nil {
			uc.logger.ErrorLevel("Error saving new user", err, map[string]interface{}{"email": email})
			return nil, err
		}

		uc.logger.InfoLevel("New user created successfully", map[string]interface{}{"user_id": domainUser.ID()})
		user = domainUser
	}

	if err = uc.accountValidation.Execute(ctx, user.Email()); err != nil {
		uc.logger.ErrorLevel("Error requesting account validation", err, map[string]interface{}{"user_id": user.ID()})
		return nil, fmt.Errorf("failed to request account validation: %w", err)
	}

	uc.logger.InfoLevel("Account validation requested successfully", map[string]interface{}{"user_id": user.ID()})
	return user, nil
}

func (uc *CreateClienteUseCase) reactivateUser(ctx context.Context, user *domain.User, name, password string, role value_objects.Role) error {
	uc.logger.InfoLevel("Reactivating user", map[string]interface{}{"user_id": user.ID(), "user_name": name})

	user.Reactivate()

	if err := user.ChangeName(name); err != nil {
		uc.logger.ErrorLevel("Error changing user name", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	if err := user.ChangeRole(role); err != nil {
		uc.logger.ErrorLevel("Error changing user role", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	newHashedPassword, err := uc.hasher.Hash(password)
	if err != nil {
		uc.logger.ErrorLevel("Error hashing new password", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	if err = user.SetPasswordHash(newHashedPassword); err != nil {
		uc.logger.ErrorLevel("Error setting new password hash", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	return uc.userRepository.Update(ctx, user)
}
