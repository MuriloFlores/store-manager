package web

import (
	"errors"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	var restErr *web_errors.RestErr

	var invalidInput *domain.ErrInvalidInput
	var notFound *domain.ErrNotFound
	var conflict *domain.ErrConflict
	var forbidden *domain.ErrForbidden
	var invalidCredentials *domain.ErrInvalidCredentials

	switch {
	case errors.As(err, &invalidInput):
		causes := []web_errors.Causes{{
			Field:   invalidInput.FieldName,
			Message: invalidInput.Reason,
		}}
		restErr = web_errors.NewBadRequestValidationError("Entrada de dados inválida", causes)

	case errors.As(err, &notFound):
		restErr = web_errors.NewNotFoundError(err.Error())

	case errors.As(err, &conflict):
		restErr = web_errors.NewConflictError(err.Error())

	case errors.As(err, &forbidden):
		restErr = web_errors.NewForbiddenError(err.Error())

	case errors.As(err, &invalidCredentials):
		restErr = web_errors.NewUnauthorizedRequestError(err.Error())

	default:
		restErr = web_errors.NewInternalServerError("Ocorreu um erro interno no servidor")
	}

	restErr.Send(w)
}
