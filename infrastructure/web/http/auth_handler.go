package http

import (
	"encoding/json"
	"github.com/muriloFlores/StoreManager/infrastructure/validation"
	"github.com/muriloFlores/StoreManager/infrastructure/web"
	"github.com/muriloFlores/StoreManager/internal/core/ports"

	dto "github.com/muriloFlores/StoreManager/infrastructure/web/DTO/authDTO"
	"github.com/muriloFlores/StoreManager/infrastructure/web/middleware"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
	"net/http"
)

type AuthHandler struct {
	useCases *auth.AuthUseCases
	logger   ports.Logger
}

func NewAuthHandler(useCases *auth.AuthUseCases, logger ports.Logger) *AuthHandler {
	return &AuthHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoLevel("Login handler invoked", nil)
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorLevel("Failed to decode request body", err)
		restErr := web_errors.NewBadRequestError("invalid body request")
		restErr.Send(w)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		h.logger.ErrorLevel("Validation error", err, map[string]interface{}{"email": req.Email, "password": req.Password})
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	token, err := h.useCases.Login.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.ErrorLevel("Login use case error", err, map[string]interface{}{"email": req.Email})
		web.HandleError(w, err)
		return
	}

	h.logger.InfoLevel("User logged in successfully", map[string]interface{}{"email": req.Email})
	respondWithJSON(w, http.StatusOK, dto.LoginResponse{Token: token})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoLevel("ChangePassword handler invoked", nil)

	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		h.logger.ErrorLevel("Failed to retrieve user identity from context", nil)
		restErr := web_errors.NewInternalServerError("unable to retrieve user identity")
		restErr.Send(w)
		return
	}

	var req dto.ChangePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorLevel("Failed to decode request body", err)
		web_errors.NewBadRequestError("invalid json body").Send(w)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		h.logger.ErrorLevel("Validation error", err, map[string]interface{}{"old_password": req.OldPassword, "new_password": req.NewPassword})
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	err := h.useCases.ChangePassword.Execute(r.Context(), actorIdentity, req.OldPassword, req.NewPassword)
	if err != nil {
		h.logger.ErrorLevel("ChangePassword use case error", err, map[string]interface{}{"user_id": actorIdentity.UserID})
		web.HandleError(w, err)
		return
	}

	h.logger.InfoLevel("Password changed successfully", map[string]interface{}{"user_id": actorIdentity.UserID})
	respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoLevel("ConfirmEmail handler invoked", nil)

	token := r.URL.Query().Get("token")

	if token == "" {
		h.logger.ErrorLevel("Token is missing in the request", nil)
		restErr := web_errors.NewBadRequestError("invalid token")
		restErr.Send(w)
		return
	}

	err := h.useCases.ConfirmUserEmailUseCase.Execute(r.Context(), token)
	if err != nil {
		h.logger.ErrorLevel("ConfirmUserEmail use case error", err, map[string]interface{}{"token": token})
		web.HandleError(w, err)
		return
	}

	h.logger.InfoLevel("Email confirmed successfully", map[string]interface{}{"token": token})

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Email verificado com sucesso!"})
}

func (h *AuthHandler) ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoLevel("ConfirmAccount handler invoked", nil)

	token := r.URL.Query().Get("token")

	if token == "" {
		h.logger.ErrorLevel("Token is missing in the request", nil)
		restErr := web_errors.NewBadRequestError("invalid token")
		restErr.Send(w)
		return
	}

	err := h.useCases.ConfirmAccountUserUseCase.Execute(r.Context(), token)
	if err != nil {
		h.logger.ErrorLevel("ConfirmAccount use case error", err, map[string]interface{}{"token": token})
		web.HandleError(w, err)
		return
	}

	h.logger.InfoLevel("Account confirmed successfully", map[string]interface{}{"token": token})

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Conta verificada com sucesso!"})
}
