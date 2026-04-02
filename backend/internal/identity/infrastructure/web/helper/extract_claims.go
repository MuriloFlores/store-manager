package helper

import (
	"context"
	"errors"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
)

var ErrUserNotInContext = errors.New("user claims not found in context")

func ExtractUserClaims(ctx context.Context) (*dto.UserClaims, error) {
	val := ctx.Value("user_claims")
	if val == nil {
		return nil, ErrUserNotInContext
	}

	claims, ok := val.(*dto.UserClaims)
	if !ok {
		return nil, ErrUserNotInContext
	}

	return claims, nil
}
