package vo

import (
	"errors"
	"strings"
)

var (
	ErrInvalidRole = errors.New("invalid role")
	ErrEmptyRole   = errors.New("empty role")
)

type Role string

const (
	EmployeeRole Role = "EMPLOYEE"
	AdminRole    Role = "ADMIN"
	ManagerRole  Role = "MANAGER"
)

func NewRole(value string) (Role, error) {
	normalizedValue := Role(strings.TrimSpace(strings.ToUpper(value)))

	if normalizedValue == "" {
		return "", ErrEmptyRole
	}

	switch normalizedValue {
	case EmployeeRole, AdminRole, ManagerRole:
		return normalizedValue, nil
	default:
		return "", ErrInvalidRole
	}
}

func (r Role) String() string {
	return string(r)
}

func AllRoles() []Role {
	return []Role{
		EmployeeRole,
		AdminRole,
		ManagerRole,
	}
}
