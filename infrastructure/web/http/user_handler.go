package http

import (
	"encoding/json"
	"github.com/muriloFlores/StoreManager/infrastructure/validation"
	"github.com/muriloFlores/StoreManager/infrastructure/web"
	dto "github.com/muriloFlores/StoreManager/infrastructure/web/DTO/userDTO"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/user"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
	"net/http"
)

type UserHandler struct {
	useCases *user.UserUseCases
}

func NewUserHandler(useCases *user.UserUseCases) *UserHandler {
	return &UserHandler{useCases: useCases}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web_errors.NewBadRequestError("invalid json body").Send(w)
		return
	}

	if err := validation.Validate.Struct(&req); err != nil {
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	createdUser, err := h.useCases.Create.Execute(r.Context(), req.Name, req.Email, req.Password, value_objects.Role(req.Role))
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := dto.UserResponse{
		ID:    createdUser.ID(),
		Name:  createdUser.Name(),
		Email: createdUser.Email(),
		Role:  createdUser.Role(),
	}

	respondWithJSON(w, http.StatusCreated, response)
}
