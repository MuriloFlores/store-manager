package product

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

var (
	ErrorIdShouldBeValid = errors.New("id should be valid")
)

type findProductByIdUseCase struct {
	productRepo repositories.ProductRepositoryInterface
}

type FindProductByIdUseCaseInterface interface {
	FindProductById(dto []DTOs.FindProductDTO) ([]DTOs.ProductDTO, error)
}

func NewFindProductByIdUseCase(productRepo repositories.ProductRepositoryInterface) FindProductByIdUseCaseInterface {
	return &findProductByIdUseCase{
		productRepo: productRepo,
	}
}

func (uc *findProductByIdUseCase) FindProductById(input []DTOs.FindProductDTO) ([]DTOs.ProductDTO, error) {
	logging.Info("FindProduct Journey", zap.String("Init", "FindProductByIdUseCase"))

	ids := make([]string, len(input))
	for i, product := range input {
		if product.Id == uuid.Nil {
			logging.Error("FindProduct Journey", zap.String("Error", "Invalid uuid"))
			return nil, ErrorIdShouldBeValid
		}

		ids[i] = product.Id.String()
	}

	productEntities, err := uc.productRepo.FindByIds(ids)
	if err != nil {
		logging.Error("FindProduct Journey", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error getting products: %w", err)
	}

	productDTOs := make([]DTOs.ProductDTO, len(productEntities))
	for i, productEntity := range productEntities {
		productDTOs[i] = DTOs.MapProductEntityToDTO(productEntity)
	}

	logging.Info("FindProduct Journey", zap.String("Finish", "FindProductByIdUseCase"))
	logging.Info("FindProduct Journey", zap.String("Finish", "FindProductByIdService"))

	return productDTOs, nil
}
