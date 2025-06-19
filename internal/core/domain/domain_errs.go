package domain

import "fmt"

type ErrInvalidInput struct {
	FieldName string
	Reason    string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input on field '%s': %s", e.FieldName, e.Reason)
}

type ErrNotFound struct {
	ResourceName string
	ResourceID   string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with identifier '%s' not found", e.ResourceName, e.ResourceID)
}

type ErrConflict struct {
	Resource string
	Details  string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("conflict on create %s: %s", e.Resource, e.Details)
}

type ErrInvalidCredentials struct{}

func (e *ErrInvalidCredentials) Error() string {
	return "invalid credentials"
}

type ErrForbidden struct {
	Action string
}

func (e *ErrForbidden) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Action)
}

type ErrInvalidToken struct {
	Reason string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token: %s", e.Reason)
}

type ErrEmailNotVerified struct{}

func (e *ErrEmailNotVerified) Error() string {
	return fmt.Sprintf("email not verified")
}
