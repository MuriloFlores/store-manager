package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/muriloFlores/StoreManager/infrastructure/validation"
	"github.com/muriloFlores/StoreManager/infrastructure/web"
	"github.com/muriloFlores/StoreManager/infrastructure/web/DTO/item_dto"
	"github.com/muriloFlores/StoreManager/infrastructure/web/DTO/pagination_dto"
	"github.com/muriloFlores/StoreManager/infrastructure/web/middleware"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/items"
	"net/http"
)

type ItemHandler struct {
	useCase *items.ItemsUseCases
	logger  ports.Logger
}

func NewItemHandler(useCase *items.ItemsUseCases, logger ports.Logger) *ItemHandler {
	return &ItemHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *ItemHandler) ListPublicItems(w http.ResponseWriter, r *http.Request) {
	params := pagination_dto.ParsePagination(r)

	paginatedResult, err := h.useCase.List.ListPublic(r.Context(), params)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	publicResponses := make([]item_dto.ClientItemResponse, 0, len(paginatedResult.Data))

	for _, itemData := range paginatedResult.Data {
		publicResponses = append(publicResponses, item_dto.ToClientItemResponse(itemData))
	}

	response := pagination_dto.PaginatedResponse[item_dto.ClientItemResponse]{
		Data:       publicResponses,
		Pagination: pagination_dto.ToPaginationInfoResponse(paginatedResult.Pagination),
	}

	respondWithJSON(w, http.StatusOK, response)

}

func (h *ItemHandler) ListInternalItems(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user identity not found in context").Send(w)
		return
	}
	params := pagination_dto.ParsePagination(r)

	paginatedResult, err := h.useCase.List.ListInternal(r.Context(), actorIdentity, params)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	internalResponse := make([]item_dto.InternalItemResponse, 0, len(paginatedResult.Data))
	for _, itemData := range paginatedResult.Data {
		internalResponse = append(internalResponse, item_dto.ToInternalItemResponse(itemData))
	}

	response := pagination_dto.PaginatedResponse[item_dto.InternalItemResponse]{
		Data:       internalResponse,
		Pagination: pagination_dto.ToPaginationInfoResponse(paginatedResult.Pagination),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user identity not found in context").Send(w)
		return
	}

	var req item_dto.CreateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web_errors.NewBadRequestError("invalid json body").Send(w)
		return
	}

	if err := validation.Validate.Struct(&req); err != nil {
		restErr := validation.TranslateError(err)
		restErr.Send(w)
		return
	}

	itemParam := items.CreateItemParams{
		Name:              req.Name,
		Description:       req.Description,
		SKU:               req.SKU,
		ItemType:          item.ItemType(req.ItemType),
		IsActive:          req.Active,
		CanBeSold:         req.CanBeSold,
		PriceSaleInCents:  req.PriceInCents,
		PriceCostInCents:  req.PriceCostInCents,
		StockQuantity:     req.StockQuantity,
		UnitOfMeasure:     req.UnitOfMeasure,
		MinimumStockLevel: req.MinimumStockLevel,
	}

	createItem, err := h.useCase.Create.Execute(r.Context(), itemParam, actorIdentity)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := item_dto.ToInternalItemResponse(createItem)

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user identity not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)
	targetID, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("item ID not provided").Send(w)
		return
	}

	if err := h.useCase.Delete.Execute(r.Context(), actorIdentity, targetID); err != nil {
		web.HandleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *ItemHandler) FindItemByID(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user identity not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)

	itemID, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("item ID not provided").Send(w)
		return
	}

	itemInfo, err := h.useCase.Find.FindByID(r.Context(), itemID, actorIdentity)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := item_dto.ToInternalItemResponse(itemInfo)

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ItemHandler) FindItemBySKU(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user identity not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)

	itemSKU, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("item SKU not provided").Send(w)
		return
	}

	itemInfo, err := h.useCase.Find.FindBySKU(r.Context(), itemSKU, actorIdentity)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := item_dto.ToInternalItemResponse(itemInfo)

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user service not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)

	targetID, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("item ID not provided").Send(w)
		return
	}

	var req item_dto.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web_errors.NewInternalServerError("invalid request body").Send(w)
		return
	}

	if err := validation.Validate.Struct(&req); err != nil {
		validation.TranslateError(err).Send(w)
		return
	}

	params := items.UpdateItemParams{
		Name:              req.Name,
		Description:       req.Description,
		IsActive:          req.Active,
		CanBeSold:         req.CanBeSold,
		PriceSaleInCents:  req.PriceSaleInCents,
		MinimumStockLevel: req.MinimumStockLevel,
	}

	fmt.Sprintf("%s", req.Active)
	fmt.Sprintf("%s", params.IsActive)

	updatedItem, err := h.useCase.Update.Execute(r.Context(), actorIdentity, targetID, params)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := item_dto.ToInternalItemResponse(updatedItem)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *ItemHandler) SearchItem(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user service not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)

	searchParam, ok := vars["param"]
	if !ok {
		web_errors.NewBadRequestError("search param not provided").Send(w)
		return
	}

	params := pagination_dto.ParsePagination(r)

	paginatedResult, err := h.useCase.Search.Execute(r.Context(), actorIdentity, searchParam, params)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	if actorIdentity != nil && actorIdentity.Role.IsStockEmployee() {
		internalResponse := make([]item_dto.InternalItemResponse, 0, len(paginatedResult.Data))
		for _, domainItem := range paginatedResult.Data {
			internalResponse = append(internalResponse, item_dto.ToInternalItemResponse(domainItem))
		}

		response := pagination_dto.PaginatedResponse[item_dto.InternalItemResponse]{
			Data:       internalResponse,
			Pagination: pagination_dto.ToPaginationInfoResponse(paginatedResult.Pagination),
		}

		respondWithJSON(w, http.StatusOK, response)
		return
	}

	publicResponse := make([]item_dto.ClientItemResponse, 0, len(paginatedResult.Data))
	for _, domainItem := range paginatedResult.Data {
		publicResponse = append(publicResponse, item_dto.ToClientItemResponse(domainItem))
	}

	response := pagination_dto.PaginatedResponse[item_dto.ClientItemResponse]{
		Data:       publicResponse,
		Pagination: pagination_dto.ToPaginationInfoResponse(paginatedResult.Pagination),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ItemHandler) ReactiveItem(w http.ResponseWriter, r *http.Request) {
	actorIdentity, ok := r.Context().Value(middleware.UserIdentityKey).(*domain.Identity)
	if !ok {
		web_errors.NewInternalServerError("user service not found in context").Send(w)
		return
	}

	vars := mux.Vars(r)

	targetID, ok := vars["id"]
	if !ok {
		web_errors.NewBadRequestError("search param not provided").Send(w)
		return
	}

	itemDomain, err := h.useCase.Reactive.Execute(r.Context(), actorIdentity, targetID)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	response := item_dto.ToInternalItemResponse(itemDomain)
	respondWithJSON(w, http.StatusOK, response)
}
