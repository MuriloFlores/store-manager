package domain

import (
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
	"time"
)

type User struct {
	id         string
	name       string
	email      string
	password   string
	role       value_objects.Role
	verifiedAt *time.Time
	deletedAt  *time.Time
}

func NewUser(id string, name, email, password string, role value_objects.Role) (*User, error) {
	if name == "" {
		return nil, &ErrInvalidInput{FieldName: "name", Reason: "name is required"}
	}
	if email == "" {
		return nil, &ErrInvalidInput{FieldName: "email", Reason: "email is required"}
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if !role.IsValid() {
		return nil, &ErrInvalidInput{FieldName: "role", Reason: "invalid role"}
	}

	return &User{
		id:         id,
		name:       name,
		email:      email,
		password:   password,
		role:       role,
		verifiedAt: nil,
	}, nil
}

func HydrateUser(id, name, email, passwordHash string, role value_objects.Role, verifiedAt, deletedAt *time.Time) *User {
	return &User{
		id:         id,
		name:       name,
		email:      email,
		password:   passwordHash,
		role:       role,
		verifiedAt: verifiedAt,
		deletedAt:  deletedAt,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Role() string {
	return u.role.ToString()
}

func (u *User) IsVerified() bool {
	return u.verifiedAt != nil
}

func (u *User) IsDeleted() bool {
	return u.deletedAt != nil
}

func (u *User) VerifiedAt() *time.Time {
	return u.verifiedAt
}

func (u *User) ChangeName(name string) error {
	if name == "" {
		return &ErrInvalidInput{FieldName: "name", Reason: "name is required"}
	}
	u.name = name
	return nil
}

func (u *User) ChangeEmail(email string) error {
	if email == "" {
		return &ErrInvalidInput{FieldName: "email", Reason: "email is required"}
	}
	u.email = email
	return nil
}

func (u *User) ChangeRole(newRole value_objects.Role) error {
	if !newRole.IsValid() {
		return &ErrInvalidInput{FieldName: "role", Reason: "invalid role"}
	}
	u.role = newRole
	return nil
}

func (u *User) SetPasswordHash(hashedPassword string) error {
	if err := validatePassword(hashedPassword); err != nil {
		return err
	}
	u.password = hashedPassword
	return nil
}

func (u *User) MarkAsVerified() {
	now := time.Now().UTC()

	u.verifiedAt = &now
}

func (u *User) Reactivate() {
	u.deletedAt = nil
}

func validatePassword(password string) error {
	if password == "" {
		return &ErrInvalidInput{FieldName: "password", Reason: "password is required"}
	}
	if len(password) < 8 {
		return &ErrInvalidInput{FieldName: "password", Reason: "password must be at least 8 characters"}
	}
	return nil
}
