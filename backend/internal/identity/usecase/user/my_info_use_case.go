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
	logger   ports.Logger
}

func NewMyInfoUseCase(userRepo ports.UserRepository, logger ports.Logger) user.MyInfoUseCase {
	return &MyInfoUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *MyInfoUseCase) Execute(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	uc.logger.Debug("fetching current user info", "userID", userID)

	userData, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to find user info", err, "userID", userID)
		return nil, err
	}

	if userData == nil {
		uc.logger.Info("user not found while fetching info", "userID", userID)
		return nil, entity.ErrUserNotFound
	}

	rolesStr := make([]string, len(userData.Roles()))
	for i, role := range userData.Roles() {
		rolesStr[i] = role.String()
	}

	uc.logger.Info("user info retrieved", "userID", userID)
	return &dto.UserInfo{
		Username: userData.Username(),
		Email:    userData.Email().String(),
		Role:     rolesStr,
	}, nil
}
