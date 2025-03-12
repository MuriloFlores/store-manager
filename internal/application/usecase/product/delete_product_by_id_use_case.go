package product

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_DTO"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

var (
	ErrorProductListNotBeEmpty = errors.New("product list not empty")
)

type deleteProductByIdUseCase struct {
	productRepo repositories.ProductRepositoryInterface
}

type DeleteProductByIdUseCaseInterface interface {
	DeleteProductById(dto []product_DTO.FindProductDTO) error
}

func NewDeleteProductByIdUseCase(productRepo repositories.ProductRepositoryInterface) DeleteProductByIdUseCaseInterface {
	return &deleteProductByIdUseCase{
		productRepo: productRepo,
	}
}

func (uc *deleteProductByIdUseCase) DeleteProductById(dto []product_DTO.FindProductDTO) error {
	logging.Info("DeleteProduct Journey", zap.String("Init", "DeleteProductByIdUseCase"))

	if len(dto) == 0 {
		return ErrorProductListNotBeEmpty
	}

	ids := make([]string, len(dto))
	for i, product := range dto {
		if product.Id == uuid.Nil {
			logging.Error("DeleteProduct Journey", zap.String("Error", "Invalid uuid"))
			return ErrorIdShouldBeValid
		}

		ids[i] = product.Id.String()
	}

	err := uc.productRepo.DeleteByIds(ids)
	if err != nil {
		logging.Error("DeleteProduct Journey", zap.String("Error", err.Error()))
		return fmt.Errorf("error deleting products: %w", err)
	}

	logging.Info("DeleteProduct Journey", zap.String("Finish", "DeleteProductByIdUseCase"))
	logging.Info("DeleteProductsByIds Journey", zap.String("Finish", "DeleteProductsByIdsService"))
	return nil
}
