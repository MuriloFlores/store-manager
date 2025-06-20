package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/muriloFlores/StoreManager/infrastructure/validation"
	"github.com/muriloFlores/StoreManager/infrastructure/web"
	dto "github.com/muriloFlores/StoreManager/infrastructure/web/DTO/userDTO"
	"github.com/muriloFlores/StoreManager/infrastructure/web/middleware"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
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

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user service not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)
	targetID, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("user ID not provided").Send(w)
		return
	}

	if err := h.useCases.Delete.Execute(r.Context(), actorIdentity, targetID); err != nil {
		web.HandleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("userInfo ID not provided").Send(w)
		return
	}

	userInfo, err := h.useCases.Find.FindByID(r.Context(), userID)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := dto.UserResponse{
		ID:    userInfo.ID(),
		Name:  userInfo.Name(),
		Email: userInfo.Email(),
		Role:  userInfo.Role(),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email, ok := vars["email"]
	if !ok {
		web_errors.NewBadRequestError("email not provided").Send(w)
	}

	userInfo, err := h.useCases.Find.FindByEmail(r.Context(), email)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := dto.UserResponse{
		ID:    userInfo.ID(),
		Name:  userInfo.Name(),
		Email: userInfo.Email(),
		Role:  userInfo.Role(),
	}

	respondWithJSON(w, http.StatusOK, response)
}
