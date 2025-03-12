package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"store-manager/internal/application/DTOs/product_DTO"
	"store-manager/internal/infrastructure/error_handler"

	"go.uber.org/zap"
	"store-manager/internal/application/services"
	"store-manager/internal/infrastructure/logging"
)

var (
	ErrorInvalidRequisition  = errors.New("invalid requisition")
	ErrorCreateProductFailed = errors.New("create product failed")
	ErrorEncodeResponseJSON  = errors.New("failed to encode response body")
	ErrorFindProductFailed   = errors.New("find product failed")
	ErrorDeletingProducts    = errors.New("deleting products failed")
	ErrorUpdateProductFailed = errors.New("update product failed")
)

type productHandler struct {
	productService services.ProductServiceInterface
}

type ProductHandlerInterface interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	FindProductById(w http.ResponseWriter, r *http.Request)
	GetAllProducts(w http.ResponseWriter, r *http.Request)
	DeleteProductsByIds(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
}

func NewProductHandler(productService services.ProductServiceInterface) ProductHandlerInterface {
	return &productHandler{productService: productService}
}

// CreateProduct godoc
// @Summary Inserir um novo produto
// @Description Insere um novo produto no sistema.
// @Tags Produtos
// @Accept json
// @Produce json
// @Param product body DTOs.CreateProductDTO true "Produto a ser inserido"
// @Success 201 {object} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /products/insert [post]
func (handler *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	logging.Info("CreateProduct Journey", zap.String("Init", "CreateProductHandler"))
	var input []product_DTO.CreateProductDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("CreateProduct Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	products, err := handler.productService.CreateProduct(input)
	if err != nil {
		logging.Error("CreateProduct Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorCreateProductFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		logging.Error("CreateProduct Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("CreateProduct Journey", zap.String("Finish", "CreateProductHandler"))
}

// FindProductById godoc
// @Summary Buscar Produtos por IDs
// @Description Recupera Produtos utilizando uma lista de IDs.
// @Tags Produtos
// @Accept json
// @Produce json
// @Param ids query string true "IDs dos Produtos (separados por vírgula)" default(8effac39-9d4d-4b20-851c-68cf0d8aae60)
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /products/get-by-ids [get]
func (handler *productHandler) FindProductById(w http.ResponseWriter, r *http.Request) {
	logging.Info("FindProductById Journey", zap.String("Init", "FindProductByIdHandler"))
	var input []product_DTO.FindProductDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("FindProductById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	products, err := handler.productService.FindProductById(input)
	if err != nil {
		logging.Error("FindProductById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorFindProductFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		logging.Error("FindProductById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("FindProductById Journey", zap.String("Finish", "FindProductByIdHandler"))
}

// GetAllProducts godoc
// @Summary Listar todos os Produtos
// @Description Retorna todos os Produtos cadastrados no sistema.
// @Tags Produtos
// @Accept json
// @Produce json
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /products/get-all [get]
func (handler *productHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	logging.Info("GetAllProducts Journey", zap.String("Init", "GetAllProductsHandler"))

	products, err := handler.productService.GetAllProducts()
	if err != nil {
		logging.Error("GetAllProducts Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorFindProductFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		logging.Error("GetAllProducts Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("GetAllProducts Journey", zap.String("Finish", "GetAllProductsHandler"))
}

// DeleteProductsByIds godoc
// @Summary Deletar Produtos por IDs
// @Description Remove vários Produtos do sistema utilizando seus IDs.
// @Tags Produtos
// @Accept json
// @Produce json
// @Param ids body []string true "Lista de IDs dos Produtos a serem deletados"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /products/delete-by-ids [delete]
func (handler *productHandler) DeleteProductsByIds(w http.ResponseWriter, r *http.Request) {
	logging.Info("DeleteProductById Journey", zap.String("Init", "DeleteProductByIdHandler"))
	var input []product_DTO.FindProductDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("DeleteProductById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	err := handler.productService.DeleteProductsByIds(input)
	if err != nil {
		logging.Error("DeleteProductById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorDeletingProducts.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]interface{}{
		"message":             "Products deleted successfully",
		"total_items_deleted": len(input),
	}

	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		logging.Error("DeleteProductById Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("DeleteProductById Journey", zap.String("Finish", "DeleteProductByIdHandler"))
}

// UpdateProduct godoc
// @Summary Atualizar Produtos
// @Description Atualiza os dados de múltiplos Produtos.
// @Tags Produtos
// @Accept json
// @Produce json
// @Param products body []DTOs.UpdateProductDTO true "Array de Produtos a serem atualizados"
// @Success 200 {array} DTOs.ProductDTO
// @Failure 400 {object} DTOs.ErrorResponse
// @Router /products/update [put]
func (handler *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	logging.Info("UpdateProducts Journey", zap.String("Init", "UpdateProductsHandler"))
	var input []product_DTO.UpdateProductDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("UpdateProducts Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusBadRequest, ErrorInvalidRequisition.Error())
		return
	}

	products, err := handler.productService.UpdateProducts(input)
	if err != nil {
		logging.Error("UpdateProducts Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorUpdateProductFailed.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		logging.Error("UpdateProducts Journey", zap.String("Error", err.Error()))
		error_handler.WriteJSONError(w, http.StatusInternalServerError, ErrorEncodeResponseJSON.Error())
		return
	}

	logging.Info("UpdateProducts Journey", zap.String("Finish", "UpdateProductsHandler"))
}
