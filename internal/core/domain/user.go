package domain

import (
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type User struct {
	id       *string
	name     string
	email    string
	password string
	role     value_objects.Role
}

func NewUser(id, name, email, password string, role value_objects.Role) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf(`name is required`)
	}
	if email == "" {
		return nil, fmt.Errorf(`email is required`)
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if !role.IsValid() {
		return nil, fmt.Errorf(`invalid role`)
	}

	return &User{
		id:       &id,
		name:     name,
		email:    email,
		password: password,
		role:     role,
	}, nil
}

func (u *User) ID() *string {
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

func (u *User) ChangeName(name string) error {
	if name == "" {
		return fmt.Errorf(`name is required`)
	}
	u.name = name
	return nil
}

func (u *User) ChangeEmail(email string) error {
	if email == "" {
		return fmt.Errorf(`email is required`)
	}
	u.email = email
	return nil
}

func (u *User) SetPasswordHash(hashedPassword string) error {
	if err := validatePassword(hashedPassword); err != nil {
		return err
	}
	u.password = hashedPassword
	return nil
}

func (u *User) ChangeRole(newRole value_objects.Role) error {
	if !newRole.IsValid() {
		return fmt.Errorf(`invalid role`)
	}
	u.role = newRole
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf(`password is required`)
	}
	if len(password) < 8 {
		return fmt.Errorf(`password must be at least 8 characters`)
	}
	return nil
}
