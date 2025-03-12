package product_raw_material_assoc

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type updateAssocUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

type UpdateAssocUseCaseInterface interface {
	UpdateByIds(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error)
}

func NewUpdateAssocUseCase(assocRepo repositories.ProductAssocRawMaterialRepositoryInterface) UpdateAssocUseCaseInterface {
	return &updateAssocUseCase{
		assocRepo: assocRepo,
	}
}

func (uc *updateAssocUseCase) UpdateByIds(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error) {
	logging.Info("UpdateAssocByIds Journey", zap.String("Init", "UpdateByIdsUseCase"))
	if len(input) == 0 {
		logging.Error("UpdateAssocById Journey", zap.String("Error", "Input can't be empty"))
		return []entity.ProductInterface{}, fmt.Errorf("invalid input")
	}

	if err := uc.assocRepo.UpdateAssociation(input); err != nil {
		logging.Error("UpdateAssocById Journey", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, fmt.Errorf("error updating associations: %s", err.Error())
	}

	productIdSet := make(map[string]struct{})
	materialIdSet := make(map[string]struct{})

	for _, assoc := range input {
		productIdSet[assoc.ProductId.String()] = struct{}{}
		materialIdSet[assoc.MaterialId.String()] = struct{}{}
	}

	var uniqueProductId, uniqueMaterialId []string

	for prodId := range productIdSet {
		uniqueProductId = append(uniqueProductId, prodId)
	}

	for materialId := range materialIdSet {
		uniqueMaterialId = append(uniqueMaterialId, materialId)
	}

	products, err := uc.assocRepo.GetAssociation(uniqueProductId, uniqueMaterialId)
	if err != nil {
		logging.Error("UpdateAssocById Journey", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, fmt.Errorf("error getting associations: %s", err.Error())
	}

	logging.Info("UpdateAssocByIds Journey", zap.String("Finish", "UpdateByIdsUseCase"))
	logging.Info("UpdateAssocByIds Journey", zap.String("Finish", "UpdateByIdsService"))
	return products, nil
}
