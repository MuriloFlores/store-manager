package product

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type updateProductUseCase struct {
	productRepo repositories.ProductRepositoryInterface
}

type UpdateProductUseCaseInterface interface {
	UpdateProduct(input []DTOs.UpdateProductDTO) ([]DTOs.ProductDTO, error)
}

func NewUpdateProductUseCase(productRepo repositories.ProductRepositoryInterface) UpdateProductUseCaseInterface {
	return &updateProductUseCase{productRepo: productRepo}
}

func (uc *updateProductUseCase) UpdateProduct(input []DTOs.UpdateProductDTO) ([]DTOs.ProductDTO, error) {
	logging.Info("UpdateProduct Journey", zap.String("Init", "UpdateProductUseCase"))
	productEntities := make([]entity.ProductInterface, len(input))

	for i, product := range input {
		id, err := entity.ParseEntityID(product.Id.String())
		if err != nil {
			return nil, ErrorIdShouldBeValid
		}

		productEntity := entity.NewProduct(
			&id,
			product.Name,
			nil,
			product.Quantity,
			product.Value.MapMoneyDTOToObject(),
		)

		productEntities[i] = productEntity
	}

	err := uc.productRepo.Update(productEntities)
	if err != nil {
		logging.Error("UpdateProduct Journey", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error while updating product: %w", err)
	}

	productDTOs := make([]DTOs.ProductDTO, len(productEntities))
	for i, product := range productEntities {
		productDTOs[i] = DTOs.MapProductEntityToDTO(product)
	}

	logging.Info("UpdateProduct Journey", zap.String("Finish", "UpdateProductUseCase"))
	logging.Info("UpdateProduct Journey", zap.String("Finish", "UpdateProductService"))

	return productDTOs, nil
}
