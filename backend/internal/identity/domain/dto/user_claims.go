package dto

import (
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Roles  []vo.Role `json:"roles"`
}
