package domain

import "fmt"

// ErrInvalidInput - Usado quando os dados de entrada falham na validação.
// Ex: um campo obrigatório está vazio, um email não tem o formato correto.
type ErrInvalidInput struct {
	FieldName string
	Reason    string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("entrada inválida no campo '%s': %s", e.FieldName, e.Reason)
}

// ErrNotFound - Usado quando um recurso específico não é encontrado.
// Ex: buscar um usuário por um ID que não existe.
type ErrNotFound struct {
	ResourceName string
	ResourceID   string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s com identificador '%s' não encontrado", e.ResourceName, e.ResourceID)
}

// ErrConflict - Usado quando uma ação viola uma regra de unicidade.
// Ex: tentar cadastrar um email que já existe.
type ErrConflict struct {
	Resource string
	Details  string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("conflito ao criar %s: %s", e.Resource, e.Details)
}

// ErrInvalidCredentials - Usado especificamente para falhas de login.
// É melhor que um ErrNotFound ou ErrInvalidInput, pois não vaza informação

type ErrInvalidCredentials struct{}

func (e *ErrInvalidCredentials) Error() string {
	return "credenciais inválidas"
}
