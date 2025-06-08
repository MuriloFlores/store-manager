package domain

import "github.com/muriloFlores/StoreManager/internal/core/value_objects"

type Identity struct {
	UserID string
	Role   value_objects.Role
}
