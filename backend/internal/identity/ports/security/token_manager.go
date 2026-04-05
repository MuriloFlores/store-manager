package security

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
)

type TokenManager interface {
	GenerateTokens(ctx context.Context, user *entity.User) (string, string, error)
	ValidateAccessToken(tokenString string) (*dto.UserClaims, error)
}
