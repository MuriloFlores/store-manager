package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/muriloFlores/StoreManager/infrastructure/validation"
	"github.com/muriloFlores/StoreManager/infrastructure/web"
	"github.com/muriloFlores/StoreManager/infrastructure/web/DTO/pagination_dto"
	dto "github.com/muriloFlores/StoreManager/infrastructure/web/DTO/user_dto"
	"github.com/muriloFlores/StoreManager/infrastructure/web/middleware"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/user"
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

	createdUser, err := h.useCases.Create.Execute(r.Context(), req.Name, req.Email, req.Password)
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

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user service not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)
	targetID, ok := vars["id"]

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web_errors.NewBadRequestError("invalid json body").Send(w)
		return
	}

	if err := validation.Validate.Struct(&req); err != nil {
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	params := user.UpdateUserParams{
		Name: req.Name,
		Role: req.Role,
	}

	responseUser, err := h.useCases.Update.Execute(r.Context(), actorIdentity, targetID, params)

	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := dto.UserResponse{
		ID:    responseUser.ID(),
		Name:  responseUser.Name(),
		Email: responseUser.Email(),
		Role:  responseUser.Role(),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) PromoteUser(w http.ResponseWriter, r *http.Request) {
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

	var req dto.PromoteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web_errors.NewBadRequestError("invalid json body").Send(w)
		return
	}

	if err := validation.Validate.Struct(&req); err != nil {
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	targetUser, err := h.useCases.Promote.Execute(r.Context(), actorIdentity, targetID, req.Role)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := dto.UserResponse{
		ID:    targetUser.ID(),
		Name:  targetUser.Name(),
		Email: targetUser.Email(),
		Role:  targetUser.Role(),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user service not found in context").Send(w)
		return
	}

	params := pagination_dto.ParsePagination(r)

	paginatedResult, err := h.useCases.List.Execute(r.Context(), actorIdentity, params)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	userResponses := make([]dto.UserResponse, 0, len(paginatedResult.Data))
	for _, userDomain := range paginatedResult.Data {
		userResponses = append(userResponses, dto.UserResponse{
			ID:    userDomain.ID(),
			Name:  userDomain.Name(),
			Email: userDomain.Email(),
			Role:  userDomain.Role(),
		})
	}

	finalResponse := pagination_dto.PaginatedResponse[dto.UserResponse]{
		Data:       userResponses,
		Pagination: pagination_dto.ToPaginationInfoResponse(paginatedResult.Pagination),
	}

	respondWithJSON(w, http.StatusOK, finalResponse)
}
