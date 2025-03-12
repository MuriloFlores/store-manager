package product_raw_material_assoc

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type findAssocByIdUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

type FindAssocByIdUseCaseInterface interface {
	FindByIds(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error)
}

func NewFindAssocByIdUseCase(assocRepo repositories.ProductAssocRawMaterialRepositoryInterface) FindAssocByIdUseCaseInterface {
	return &findAssocByIdUseCase{
		assocRepo: assocRepo,
	}
}

func (uc *findAssocByIdUseCase) FindByIds(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error) {
	logging.Info("FindAssocById Journey", zap.String("Init", "FindByIdAssocUseCase"))
	if len(input) == 0 {
		logging.Error("FindAssocById Journey", zap.String("Error", "Input can't be empty"))
		return []entity.ProductInterface{}, nil
	}

	var productId, materialId []string

	for _, assoc := range input {
		if assoc.MaterialId == uuid.Nil || assoc.ProductId == uuid.Nil {
			logging.Error("FindAssocById Journey", zap.String("Error", "Invalid ids"))
			return []entity.ProductInterface{}, errors.New("invalid ids")
		}

		productId = append(productId, assoc.ProductId.String())
		materialId = append(materialId, assoc.MaterialId.String())
	}

	products, err := uc.assocRepo.GetAssociation(productId, materialId)
	if err != nil {
		logging.Error("FindAssocById Journey", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, errors.New("error finding products")
	}

	logging.Info("FindAssocById Journey", zap.String("Finish", "FindByIdAssocUseCase"))
	logging.Info("FindAssocById Journey", zap.String("Finish", "FindByIdAssocService"))
	return products, nil
}
