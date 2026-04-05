package model

import (
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/google/uuid"
)

type UserModel struct {
	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	Email          string     `bun:"email,notnull,unique"`
	Username       string     `bun:"username,notnull"`
	Password       string     `bun:"password,notnull"`
	Roles          []string   `bun:"roles,array,notnull"`
	Active         bool       `bun:"active,notnull"`
	FailedAttempts int        `bun:"failed_attempts,notnull"`
	LockedUntil    *time.Time `bun:"locked_until"`
	EmailVerified  bool       `bun:"email_verified,notnull"`
}

func ToModel(u *entity.User) *UserModel {
	rolesStr := make([]string, 0, len(u.Roles()))
	for _, role := range u.Roles() {
		rolesStr = append(rolesStr, role.String())
	}

	return &UserModel{
		ID:             u.ID(),
		Email:          u.Email().String(),
		Username:       u.Username(),
		Password:       u.Password().String(),
		Roles:          rolesStr,
		Active:         u.IsActive(),
		FailedAttempts: u.FailedAttempts(),
		LockedUntil:    u.LockedUntil(),
		EmailVerified:  u.EmailVerified(),
	}
}

func ToEntity(m *UserModel) (*entity.User, error) {
	return entity.RestoreUser(
		m.ID,
		m.Email,
		m.Username,
		m.Password,
		m.Roles,
		m.Active,
		m.FailedAttempts,
		m.LockedUntil,
		m.EmailVerified,
	)
}
