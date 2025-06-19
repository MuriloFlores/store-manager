package http

import (
	"encoding/json"
	"github.com/muriloFlores/StoreManager/infrastructure/validation"
	"github.com/muriloFlores/StoreManager/infrastructure/web"

	dto "github.com/muriloFlores/StoreManager/infrastructure/web/DTO/authDTO"
	"github.com/muriloFlores/StoreManager/infrastructure/web/middleware"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
	"net/http"
)

type AuthHandler struct {
	useCases *auth.AuthUseCases
}

func NewAuthHandler(useCases *auth.AuthUseCases) *AuthHandler {
	return &AuthHandler{
		useCases: useCases,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		restErr := web_errors.NewBadRequestError("invalid body request")
		restErr.Send(w)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	token, err := h.useCases.Login.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, dto.LoginResponse{Token: token})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		restErr := web_errors.NewInternalServerError("unable to retrieve user identity")
		restErr.Send(w)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web_errors.NewBadRequestError("invalid json body").Send(w)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	err := h.useCases.ChangePassword.Execute(r.Context(), actorIdentity, req.OldPassword, req.NewPassword)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
