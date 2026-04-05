package entity

import (
	"errors"
	"slices"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
)

var (
	ErrEmptyUsername = errors.New("empty username")
)

type User struct {
	id             uuid.UUID
	email          vo.Email
	username       string
	password       vo.Password
	roles          []vo.Role
	active         bool
	failedAttempts int
	lockedUntil    *time.Time
	emailVerified  bool
}

func NewUser(email vo.Email, username string, password vo.Password, roles []vo.Role) (*User, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	return &User{
		id:             uuid.New(),
		email:          email,
		username:       username,
		password:       password,
		roles:          roles,
		active:         true,
		failedAttempts: 0,
		lockedUntil:    nil,
		emailVerified:  false,
	}, nil
}

func RestoreUser(
	id uuid.UUID,
	email string,
	username string,
	password string,
	roles []string,
	active bool,
	failedAttempts int,
	lockedUntil *time.Time,
	emailVerified bool,
) (*User, error) {
	restoredPassword, err := vo.RestorePassword(password)
	if err != nil {
		return nil, err
	}

	restoredRoles := make([]vo.Role, 0, len(roles))
	for _, role := range roles {
		restRole, err := vo.NewRole(role)
		if err != nil {
			return nil, err
		}

		restoredRoles = append(restoredRoles, restRole)
	}

	restEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:             id,
		email:          restEmail,
		username:       username,
		password:       restoredPassword,
		roles:          restoredRoles,
		active:         active,
		failedAttempts: failedAttempts,
		lockedUntil:    lockedUntil,
		emailVerified:  emailVerified,
	}, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() vo.Email {
	return u.email
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() vo.Password {
	return u.password
}

func (u *User) Roles() []vo.Role {
	return u.roles
}

func (u *User) Activate() {
	u.active = true
}

func (u *User) FailedAttempts() int {
	return u.failedAttempts
}

func (u *User) LockedUntil() *time.Time {
	return u.lockedUntil
}

func (u *User) EmailVerified() bool {
	return u.emailVerified
}

func (u *User) Deactivate() {
	u.active = false
}

func (u *User) IsActive() bool {
	return u.active
}

func (u *User) ChangeEmail(email vo.Email) {
	u.email = email
}

func (u *User) ChangePassword(p vo.Password) {
	u.password = p
}

func (u *User) AddRole(r vo.Role) {
	for _, role := range u.roles {
		if role == r {
			return
		}
	}

	u.roles = append(u.roles, r)
}

func (u *User) RemoveRole(r vo.Role) {
	u.roles = slices.DeleteFunc(u.roles, func(role vo.Role) bool {
		return role == r
	})
}

func (u *User) ReplaceRoles(newRoles []vo.Role) {
	if len(newRoles) == 0 {
		u.roles = []vo.Role{vo.EmployeeRole}
		return
	}

	u.roles = newRoles
}

func (u *User) IsLocked(now time.Time) bool {
	if u.lockedUntil == nil {
		return false
	}

	return u.lockedUntil.After(now)
}

func (u *User) RecordFailedLogin(threshold int, baseDuration time.Duration, now time.Time) {
	u.failedAttempts++

	if u.failedAttempts >= threshold {
		fator := u.failedAttempts / threshold
		extraDuration := baseDuration * time.Duration(fator)

		expiration := now.Add(extraDuration)
		u.lockedUntil = &expiration
	}
}

func (u *User) ResetFailedAttempts() {
	u.failedAttempts = 0
	u.lockedUntil = nil
}

func (u *User) VerifyEmail() {
	u.emailVerified = true
}
