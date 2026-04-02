package helper

import (
	"errors"
	"net/http"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch {
	// 1. Entidade Improcessável (422) - Erros de Validação e Regras de Domínio
	case errors.Is(err, vo.ErrPasswordTooShort),
		errors.Is(err, vo.ErrLowPasswordComplexity),
		errors.Is(err, vo.ErrEmptyPassword),
		errors.Is(err, vo.ErrInvalidRole),
		errors.Is(err, vo.ErrEmptyRole),
		errors.Is(err, vo.ErrEmptyEmail),
		errors.Is(err, vo.ErrInvalidEmail),
		errors.Is(err, vo.ErrInvalidOTPFormat),
		errors.Is(err, entity.ErrEmptyUsername):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

	// 2. Não Autorizado (401) - Falhas de Autenticação e Tokens
	case errors.Is(err, entity.ErrInvalidCredentials),
		errors.Is(err, entity.ErrInvalidOldPassword),
		errors.Is(err, entity.ErrUserIsDeactivated),
		errors.Is(err, infrastructure.ErrInvalidToken),
		errors.Is(err, infrastructure.ErrExpiredToken),
		errors.Is(err, infrastructure.ErrUnexpectedMethod):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	// 3. Não Encontrado (404)
	case errors.Is(err, entity.ErrUserNotFound),
		errors.Is(err, entity.ErrSessionNotFound),
		errors.Is(err, entity.ErrOTPNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	// 4. Catch-all para Erros Internos (500)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error, please try again later"})
	}
}
