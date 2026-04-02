package user

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/user"
	"github.com/google/uuid"
)

type MyInfoUseCase struct {
	userRepo ports.UserRepository
}

func NewMyInfoUseCase(userRepo ports.UserRepository) user.MyInfoUseCase {
	return &MyInfoUseCase{
		userRepo: userRepo,
	}
}

func (uc *MyInfoUseCase) Execute(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	userData, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userData == nil {
		return nil, entity.ErrUserNotFound
	}

	rolesStr := make([]string, len(userData.Roles()))
	for i, role := range userData.Roles() {
		rolesStr[i] = role.String()
	}

	return &dto.UserInfo{
		Username: userData.Username(),
		Email:    userData.Email().String(),
		Role:     rolesStr,
	}, nil
}
