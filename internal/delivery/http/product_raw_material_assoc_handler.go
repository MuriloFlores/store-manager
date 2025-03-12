package http

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"store-manager/internal/application/DTOs/product_DTO"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	"store-manager/internal/application/services"
	"store-manager/internal/infrastructure/error_handler"
	"store-manager/internal/infrastructure/logging"
)

var (
	ErrorCreateAssocFailed = errors.New("falha ao criar associação")
	ErrorFindAssocFailed   = errors.New("falha ao buscar associação")
	ErrorDeletingAssoc     = errors.New("falha ao deletar associação")
	ErrorUpdateAssocFailed = errors.New("falha ao atualizar associação")
)

type assocHandler struct {
	assocService services.ProductRawMaterialAssocServiceInterface
}

type AssocHandlerInterface interface {
	CreateAssoc(w http.ResponseWriter, r *http.Request)
	DeleteAssoc(w http.ResponseWriter, r *http.Request)
	FindAllAssociations(w http.ResponseWriter, r *http.Request)
	FindByCriteria(w http.ResponseWriter, r *http.Request)
	FindAssocById(w http.ResponseWriter, r *http.Request)
	UpdateAssoc(w http.ResponseWriter, r *http.Request)
}

func NewProductRawMaterialAssocHandler(assocService services.ProductRawMaterialAssocServiceInterface) AssocHandlerInterface {
	return &assocHandler{
		assocService: assocService,
	}
}

// CreateAssoc godoc
// @Summary Criar Associações
// @Description Cria associações entre produtos e matérias-primas.
// @Tags Associações
// @Accept json
// @Produce json
// @Param associations body []DTOs.ProductAssocRawMaterialDTO true "Array de associações a serem criadas"
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Failure 500 {object} DTOs.ErrorResponse
// @Router /assoc-product-material/insert [post]
func (handler *assocHandler) CreateAssoc(w http.ResponseWriter, r *http.Request) {
	logging.Info("CreateAssoc Journey", zap.String("Init", "CreateAssocHandler"))
	var input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("CreateAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	productsEntities, err := handler.assocService.CreateAssoc(input)
	if err != nil {
		logging.Error("CreateAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorCreateAssocFailed.Error())
		return
	}

	productsDTOs := make([]product_DTO.ProductDTO, len(productsEntities))
	for i, product := range productsEntities {
		productsDTOs[i] = product_DTO.MapProductEntityToDTO(product)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productsDTOs); err != nil {
		logging.Error("CreateAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorCreateAssocFailed.Error())
		return
	}

	logging.Info("CreateAssoc Journey", zap.String("Finish", "CreateAssocHandler"))
}

// FindAllAssociations godoc
// @Summary Obter Todas as Associações
// @Description Recupera todas as associações entre produtos e matérias-primas e retorna os produtos enriquecidos.
// @Tags Associações
// @Produce json
// @Success 200 {array} DTOs.ProductDTO
// @Failure 500 {object} DTOs.ErrorResponse
// @Router /assoc-product-material/get-all [get]
func (handler *assocHandler) FindAllAssociations(w http.ResponseWriter, r *http.Request) {
	logging.Info("FindAllAssociations Journey", zap.String("Init", "FindAllAssociationsHandler"))

	productsEntities, err := handler.assocService.FindAllAssociations()
	if err != nil {
		logging.Error("FindAllAssociations Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorFindAssocFailed.Error())
		return
	}

	productsDTOs := make([]product_DTO.ProductDTO, len(productsEntities))
	for i, product := range productsEntities {
		productsDTOs[i] = product_DTO.MapProductEntityToDTO(product)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productsDTOs); err != nil {
		logging.Error("FindAllAssociations Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorFindAssocFailed.Error())
		return
	}

	logging.Info("FindAllAssociations Journey", zap.String("Finish", "FindAllAssociationsHandler"))
}

// FindByCriteria godoc
// @Summary Buscar Associações por Critério
// @Description Recupera associações com base nos critérios e informações fornecidos.
// @Tags Associações
// @Accept json
// @Produce json
// @Param criteria body DTOs.FindAssocByCriteriaDTO true "Critérios para buscar associações"
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Failure 500 {object} DTOs.ErrorResponse
// @Router /assoc-product-material/get-by-criteria [get]
func (handler *assocHandler) FindByCriteria(w http.ResponseWriter, r *http.Request) {
	logging.Info("FindByCriteria Journey", zap.String("Init", "FindByCriteriaHandler"))
	var input product_assoc_raw_material_DTO.FindAssocByCriteriaDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("FindByCriteria Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	productsEntities, err := handler.assocService.FindByCriteria(input.Criteria, input.Info)
	if err != nil {
		logging.Error("FindByCriteria Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorFindAssocFailed.Error())
		return
	}

	productsDTOs := make([]product_DTO.ProductDTO, len(productsEntities))
	for i, product := range productsEntities {
		productsDTOs[i] = product_DTO.MapProductEntityToDTO(product)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productsDTOs); err != nil {
		logging.Error("FindByCriteria Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("FindByCriteria Journey", zap.String("Finish", "FindByCriteriaHandler"))
}

// FindAssocById godoc
// @Summary Buscar Associação por ID
// @Description Recupera associações com base no DTO fornecido.
// @Tags Associações
// @Accept json
// @Produce json
// @Param association body []DTOs.ProductAssocRawMaterialDTO true "DTO de associação para buscar associações"
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Failure 500 {object} DTOs.ErrorResponse
// @Router /assoc-product-material/get-by-ids [get]
func (handler *assocHandler) FindAssocById(w http.ResponseWriter, r *http.Request) {
	logging.Info("FindAssocById Journey", zap.String("Init", "FindAssocByIdHandler"))
	var input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("FindAssocById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	productEntities, err := handler.assocService.FindAssocById(input)
	if err != nil {
		logging.Error("FindAssocById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorFindAssocFailed.Error())
		return
	}

	productsDTOs := make([]product_DTO.ProductDTO, len(productEntities))
	for i, product := range productEntities {
		productsDTOs[i] = product_DTO.MapProductEntityToDTO(product)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productsDTOs); err != nil {
		logging.Error("FindAssocById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("FindAssocById Journey", zap.String("Finish", "FindAssocByIdHandler"))
}

// DeleteAssoc godoc
// @Summary Deletar Associações
// @Description Deleta associações com base nos IDs de produtos e matérias-primas.
// @Tags Associações
// @Accept json
// @Produce json
// @Param criteria body DTOs.FindProductAssocRawMaterialDTO true "Critérios para deletar associações"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} DTOs.ErrorResponse
// @Failure 500 {object} DTOs.ErrorResponse
// @Router /assoc-product-material/delete-by-ids [delete]
func (handler *assocHandler) DeleteAssoc(w http.ResponseWriter, r *http.Request) {
	logging.Info("DeleteAssoc Journey", zap.String("Init", "DeleteAssocHandler"))
	var input product_assoc_raw_material_DTO.FindProductAssocRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("DeleteAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	productIds := make([]string, len(input.ProductIds))
	materialIds := make([]string, len(input.MaterialIds))

	for i, id := range input.ProductIds {
		productIds[i] = id.String()
	}

	for i, id := range input.MaterialIds {
		materialIds[i] = id.String()
	}

	err := handler.assocService.DeleteAssoc(productIds, materialIds)
	if err != nil {
		logging.Error("DeleteAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorDeletingAssoc.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var deletedItems int

	if len(input.ProductIds) > len(input.MaterialIds) {
		deletedItems = len(input.ProductIds)
	} else {
		deletedItems = len(input.MaterialIds)
	}

	responseMessage := map[string]interface{}{
		"message":             "Associações deletadas com sucesso",
		"total_items_deleted": deletedItems,
	}

	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		logging.Error("DeleteAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("DeleteAssoc Journey", zap.String("Finish", "DeleteAssocHandler"))
}

// UpdateAssoc godoc
// @Summary Atualizar Associações
// @Description Atualiza as associações e retorna os produtos atualizados com as matérias-primas associadas.
// @Tags Associações
// @Accept json
// @Produce json
// @Param associations body []DTOs.ProductAssocRawMaterialDTO true "Array de associações para atualizar"
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Failure 500 {object} DTOs.ErrorResponse
// @Router /assoc-product-material/update [put]
func (handler *assocHandler) UpdateAssoc(w http.ResponseWriter, r *http.Request) {
	logging.Info("UpdateAssoc Journey", zap.String("Init", "UpdateAssocHandler"))
	var input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("UpdateAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	productEntities, err := handler.assocService.UpdateAssoc(input)
	if err != nil {
		logging.Error("UpdateAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorUpdateAssocFailed.Error())
		return
	}

	productsDTOs := make([]product_DTO.ProductDTO, len(productEntities))
	for i, product := range productEntities {
		productsDTOs[i] = product_DTO.MapProductEntityToDTO(product)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productsDTOs); err != nil {
		logging.Error("UpdateAssoc Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("UpdateAssoc Journey", zap.String("Finish", "UpdateAssocHandler"))
}
