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
	userRepo ports.UserRepository
	pepper   string
}

func NewCreateUserService(userRepo ports.UserRepository, pepper string) user.CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
		pepper:   pepper,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) error {
	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return err
	}

	password, err := vo.NewPassword(input.Password, uc.pepper)
	if err != nil {
		return err
	}

	roles := make([]vo.Role, 0, len(input.Roles))
	for _, role := range input.Roles {
		voRole, err := vo.NewRole(role)
		if err != nil {
			return err
		}

		roles = append(roles, voRole)
	}

	user, err := entity.NewUser(
		email,
		input.Username,
		password,
		roles,
	)
	if err != nil {
		return err
	}

	return uc.userRepo.Save(ctx, user)
}
