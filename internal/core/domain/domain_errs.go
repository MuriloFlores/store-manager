package domain

import "fmt"

// ErrInvalidInput - Usado quando os dados de entrada falham na validação.
// Ex: um campo obrigatório está vazio, um email não tem o formato correto.
type ErrInvalidInput struct {
	FieldName string
	Reason    string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input on field '%s': %s", e.FieldName, e.Reason)
}

// ErrNotFound - Usado quando um recurso específico não é encontrado.
// Ex: buscar um usuário por um ID que não existe.
type ErrNotFound struct {
	ResourceName string
	ResourceID   string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with identifier '%s' not found", e.ResourceName, e.ResourceID)
}

// ErrConflict - Usado quando uma ação viola uma regra de unicidade.
// Ex: tentar cadastrar um email que já existe.
type ErrConflict struct {
	Resource string
	Details  string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("conflict on create %s: %s", e.Resource, e.Details)
}

// ErrInvalidCredentials - Usado especificamente para falhas de login.
// É melhor que um ErrNotFound ou ErrInvalidInput, pois não vaza informação

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

// ErrInvalidToken é usado quando um token (de reset, de confirmação, etc.) é inválido,
// não encontrado ou expirado.
type ErrInvalidToken struct {
	Reason string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("token inválido: %s", e.Reason)
}
