package product

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

var (
	ErrorNameIsRequired      = errors.New("name is required")
	ErrorValueMustBePositive = errors.New("value must be positive")
)

type createProductUseCase struct {
	productRepo repositories.ProductRepositoryInterface
}

type CreateProductUseCaseInterface interface {
	CreateProduct(input []product_DTO.CreateProductDTO) ([]product_DTO.ProductDTO, error)
}

func NewCreateProductUseCase(productRepo repositories.ProductRepositoryInterface) CreateProductUseCaseInterface {
	return &createProductUseCase{productRepo: productRepo}
}

func (uc *createProductUseCase) CreateProduct(input []product_DTO.CreateProductDTO) ([]product_DTO.ProductDTO, error) {
	logging.Info("CreateProduct Journey", zap.String("Init", "CreateProductUseCase"))
	productEntities := make([]entity.ProductInterface, len(input))
	for i, product := range input {
		if product.Name == "" {
			return []product_DTO.ProductDTO{}, ErrorNameIsRequired
		}

		if product.Value.TotalCents <= 0 {
			return []product_DTO.ProductDTO{}, ErrorValueMustBePositive
		}

		productEntity := entity.NewProduct(
			nil,
			product.Name,
			[]entity.RawMaterialInterface{},
			product.Quantity,
			product.Value.MapMoneyDTOToObject(),
		)

		productEntities[i] = productEntity
	}

	err := uc.productRepo.Save(productEntities)
	if err != nil {
		return []product_DTO.ProductDTO{}, fmt.Errorf("error saving product: %w", err)
	}

	productDTO := make([]product_DTO.ProductDTO, len(productEntities))

	for i, product := range productEntities {
		productDTO[i] = product_DTO.MapProductEntityToDTO(product)
	}

	logging.Info("CreateProduct Journey", zap.String("Finish", "CreateProductUseCase"))
	logging.Info("CreateProduct Journey", zap.String("Finish", "CreateProductService"))

	return productDTO, nil
}
