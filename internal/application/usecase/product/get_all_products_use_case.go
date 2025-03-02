package product

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type getAllProductsUseCase struct {
	productRepo repositories.ProductRepositoryInterface
}

type GetAllProductsUseCase interface {
	GetAllProducts() ([]DTOs.ProductDTO, error)
}

func NewGetAllProducts(productRepo repositories.ProductRepositoryInterface) GetAllProductsUseCase {
	return &getAllProductsUseCase{
		productRepo: productRepo,
	}
}

func (uc *getAllProductsUseCase) GetAllProducts() ([]DTOs.ProductDTO, error) {
	logging.Info("GetAllProducts Journey", zap.String("Init", "GetAllProductsUseCase"))

	productEntities, err := uc.productRepo.GetAllProducts()
	if err != nil {
		logging.Error("GetAllProducts Journey", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error while getting all products: %w", err)
	}

	productDTOs := make([]DTOs.ProductDTO, len(productEntities))
	for i, productEntity := range productEntities {
		productDTOs[i] = DTOs.MapProductEntityToDTO(productEntity)
	}

	logging.Info("GetAllProducts Journey", zap.String("Finish", "GetAllProductsUseCase"))
	logging.Info("GetAllProducts Journey", zap.String("Finish", "GetAllProductsService"))
	return productDTOs, nil
}
