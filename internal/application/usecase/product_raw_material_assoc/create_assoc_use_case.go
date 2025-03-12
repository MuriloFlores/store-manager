package product_raw_material_assoc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type createAssocUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

type CreateAssocUseCaseInterface interface {
	CreateAssoc(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error)
}

func NewCreateProductRawMaterialAssocUseCase(assocRepo repositories.ProductAssocRawMaterialRepositoryInterface) CreateAssocUseCaseInterface {
	return &createAssocUseCase{
		assocRepo: assocRepo,
	}
}

func (uc *createAssocUseCase) CreateAssoc(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error) {
	logging.Info("CreateAssoc Journey", zap.String("Init", "CreateAssocUseCase"))
	if len(input) == 0 {
		logging.Error("CreateAssoc Journey", zap.String("Error", "input can't be empty"))
		return []entity.ProductInterface{}, fmt.Errorf("input is empty")
	}

	for i, assoc := range input {
		if assoc.ProductId == uuid.Nil {
			return []entity.ProductInterface{}, fmt.Errorf("association at index %d: productId is empty", i)
		}

		if assoc.MaterialId == uuid.Nil {
			return []entity.ProductInterface{}, fmt.Errorf("association at index %d: materialId is empty", i)
		}

		if assoc.QuantityUsed <= 0 {
			return []entity.ProductInterface{}, fmt.Errorf("association at index %d: quantity is empty", i)
		}
	}

	if err := uc.assocRepo.CreateAssociation(input); err != nil {
		logging.Error("CreateAssoc Journey", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, err
	}

	productIdSet := make(map[string]struct{})
	materialIdSet := make(map[string]struct{})

	for _, assoc := range input {
		productIdSet[assoc.ProductId.String()] = struct{}{}
		materialIdSet[assoc.MaterialId.String()] = struct{}{}
	}

	var uniqueProductIds, uniqueMaterialIds []string

	for productId := range productIdSet {
		uniqueProductIds = append(uniqueProductIds, productId)
	}

	for materialId := range materialIdSet {
		uniqueMaterialIds = append(uniqueMaterialIds, materialId)
	}

	products, err := uc.assocRepo.GetAssociation(uniqueProductIds, uniqueMaterialIds)
	if err != nil {
		logging.Error("CreateAssoc Journey", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, errors.Errorf("error retrieving associations: %v", err)
	}

	logging.Info("CreateAssoc Journey", zap.String("Finish", "CreateAssocUseCase"))
	logging.Info("CreateAssoc Journey", zap.String("Finish", "CreateAssocService"))
	return products, nil
}
