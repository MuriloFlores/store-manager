package user

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/user"
)

type CreateUserUseCase struct {
	userRepo  ports.UserRepository
	logger    ports.Logger
	txManager ports.TransactionManager
	pepper    string
}

func NewCreateUserService(userRepo ports.UserRepository, logger ports.Logger, txManager ports.TransactionManager, pepper string) user.CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:  userRepo,
		logger:    logger,
		txManager: txManager,
		pepper:    pepper,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) error {
	uc.logger.Debug("starting user creation", "username", input.Username, "email", input.Email)

	email, err := vo.NewEmail(input.Email)
	if err != nil {
		uc.logger.Info("invalid email format in creation", "email", input.Email)
		return err
	}

	password, err := vo.NewPassword(input.Password, uc.pepper)
	if err != nil {
		uc.logger.Error("failed to process password", err, "email", input.Email)
		return err
	}

	roles := make([]vo.Role, 0, len(input.Roles))
	for _, role := range input.Roles {
		voRole, err := vo.NewRole(role)
		if err != nil {
			uc.logger.Info("invalid role in user creation", "role", role)
			return err
		}

		roles = append(roles, voRole)
	}

	createdUser, err := entity.NewUser(
		email,
		input.Username,
		password,
		roles,
	)
	if err != nil {
		uc.logger.Error("failed to create user entity", err, "email", input.Email)
		return err
	}

	err = uc.txManager.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.userRepo.Save(txCtx, createdUser); err != nil {
			uc.logger.Error("failed to save user in repository", err, "email", input.Email)
			return err
		}

		return nil
	})

	uc.logger.Info("user created successfully", "userID", createdUser.ID(), "email", input.Email)
	return nil
}
