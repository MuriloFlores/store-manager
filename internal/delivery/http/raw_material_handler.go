package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"store-manager/internal/application/DTOs"
	"store-manager/internal/application/services"
	"store-manager/internal/infrastructure/error_handler"
	"store-manager/internal/infrastructure/logging"

	"go.uber.org/zap"
)

var (
	ErrorCreateRawMaterialsFailed = errors.New("create raw material failed")
	ErrorFindRawMaterialsFailed   = errors.New("find raw material failed")
	ErrorDeletingRawMaterials     = errors.New("deleting raw material failed")
	ErrorUpdateRawMaterialsFailed = errors.New("update raw material failed")
)

type rawMaterialHandler struct {
	rawMaterialService services.RawMaterialServiceInterface
}

type RawMaterialHandlerInterface interface {
	CreateRawMaterial(w http.ResponseWriter, r *http.Request)
	FindRawMaterial(w http.ResponseWriter, r *http.Request)
	GetAllRawMaterials(w http.ResponseWriter, r *http.Request)
	DeleteRawMaterial(w http.ResponseWriter, r *http.Request)
	UpdateRawMaterial(w http.ResponseWriter, r *http.Request)
}

func NewRawMaterialHandler(rawMaterialService services.RawMaterialServiceInterface) RawMaterialHandlerInterface {
	return &rawMaterialHandler{
		rawMaterialService: rawMaterialService,
	}
}

// CreateRawMaterial godoc
// @Summary Inserir novas matérias-primas
// @Description Insere uma ou mais matérias-primas no sistema.
// @Tags matéria-prima
// @Accept json
// @Produce json
// @Param rawMaterials body []DTOs.CreateRawMaterialDTO true "Array de matérias-primas a serem inseridas"
// @Success 201 {array} DTOs.RawMaterialDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /raw-material/insert [post]
func (handler *rawMaterialHandler) CreateRawMaterial(w http.ResponseWriter, r *http.Request) {
	logging.Info("CreateRawMaterial Journey", zap.String("Init", "CreateRawMaterialHandler"))
	var input []DTOs.CreateRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("CreateRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	rawMaterials, err := handler.rawMaterialService.CreateRawMaterial(input)
	if err != nil {
		logging.Error("CreateRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorCreateRawMaterialsFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(rawMaterials); err != nil {
		logging.Error("CreateRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("CreateRawMaterial Journey", zap.String("Finish", "CreateRawMaterialHandler"))
}

// FindRawMaterial godoc
// @Summary Buscar matérias-primas por IDs
// @Description Recupera matérias-primas utilizando uma lista de IDs.
// @Tags matéria-prima
// @Accept json
// @Produce json
// @Param ids query string true "IDs das matérias-primas (separados por vírgula)" default(8effac39-9d4d-4b20-851c-68cf0d8aae60)
// @Success 200 {array} DTOs.RawMaterialDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /raw-material/get-by-ids [get]
func (handler *rawMaterialHandler) FindRawMaterial(w http.ResponseWriter, r *http.Request) {
	logging.Info("FindRawMaterial Journey", zap.String("Init", "FindRawMaterialHandler"))
	var input []DTOs.FindRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("FindRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	rawMaterials, err := handler.rawMaterialService.FindRawMaterialById(input)
	if err != nil {
		logging.Error("FindRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorFindRawMaterialsFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rawMaterials); err != nil {
		logging.Error("FindRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("FindRawMaterial Journey", zap.String("Finish", "FindRawMaterialHandler"))
}

// GetAllRawMaterials godoc
// @Summary Listar todas as matérias-primas
// @Description Retorna todas as matérias-primas cadastradas no sistema.
// @Tags matéria-prima
// @Accept json
// @Produce json
// @Success 200 {array} DTOs.RawMaterialDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /raw-material/get-all [get]
func (handler *rawMaterialHandler) GetAllRawMaterials(w http.ResponseWriter, r *http.Request) {
	logging.Info("GetAllRawMaterials Journey", zap.String("Init", "GetAllRawMaterials"))

	rawMaterials, err := handler.rawMaterialService.GetAllRawMaterials()
	if err != nil {
		logging.Error("GetAllRawMaterials Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorFindRawMaterialsFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rawMaterials); err != nil {
		logging.Error("GetAllRawMaterials Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("GetAllRawMaterials Journey", zap.String("Finish", "GetAllRawMaterials"))
}

// DeleteRawMaterial godoc
// @Summary Deletar matérias-primas por IDs
// @Description Remove várias matérias-primas do sistema utilizando seus IDs.
// @Tags matéria-prima
// @Accept json
// @Produce json
// @Param input body []DTOs.FindRawMaterialDTO true "Lista de IDs das matérias-primas a serem deletadas"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /raw-material/delete-by-ids [delete]
func (handler *rawMaterialHandler) DeleteRawMaterial(w http.ResponseWriter, r *http.Request) {
	logging.Info("DeleteRawMaterial Journey", zap.String("Init", "DeleteRawMaterial"))
	var input []DTOs.FindRawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("DeleteRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	err := handler.rawMaterialService.DeleteRawMaterial(input)
	if err != nil {
		logging.Error("DeleteRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorDeletingRawMaterials.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]interface{}{
		"message":             "RawMaterials Deleted Successfully",
		"total_items_deleted": len(input),
	}

	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		logging.Error("DeleteRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("DeleteRawMaterial Journey", zap.String("Finish", "DeleteRawMaterialHandler"))
}

// UpdateRawMaterial godoc
// @Summary Atualizar matérias-primas
// @Description Atualiza os dados de múltiplas matérias-primas.
// @Tags matéria-prima
// @Accept json
// @Produce json
// @Param rawMaterials body []DTOs.RawMaterialDTO true "Array de matérias-primas a serem atualizadas"
// @Success 200 {array} DTOs.RawMaterialDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /raw-material/update [put]
func (handler *rawMaterialHandler) UpdateRawMaterial(w http.ResponseWriter, r *http.Request) {
	logging.Info("UpdateRawMaterial Journey", zap.String("Init", "UpdateRawMaterial"))
	var input []DTOs.RawMaterialDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("UpdateRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	rawMaterials, err := handler.rawMaterialService.UpdateRawMaterial(input)
	if err != nil {
		logging.Error("UpdateRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorUpdateRawMaterialsFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rawMaterials); err != nil {
		logging.Error("UpdateRawMaterial Journey", zap.Error(err))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("UpdateRawMaterial Journey", zap.String("Finish", "UpdateRawMaterialHandler"))
}
