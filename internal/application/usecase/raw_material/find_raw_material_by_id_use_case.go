package raw_material

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type findRawMaterialByIdUseCase struct {
	productRepo repositories.RawMaterialsRepositoryInterface
}

type FindRawMaterialByIdUseCaseInterface interface {
	FindRawMaterialById(input []DTOs.FindRawMaterialDTO) ([]DTOs.RawMaterialDTO, error)
}

func NewFindRawMaterialByIdUseCase(productRepo repositories.RawMaterialsRepositoryInterface) FindRawMaterialByIdUseCaseInterface {
	return &findRawMaterialByIdUseCase{
		productRepo: productRepo,
	}
}

func (uc *findRawMaterialByIdUseCase) FindRawMaterialById(input []DTOs.FindRawMaterialDTO) ([]DTOs.RawMaterialDTO, error) {
	logging.Info("FindRawMaterialById Journey", zap.String("Init", "FindRawMaterialByIdUseCase"))

	ids := make([]string, len(input))
	for i, rawMaterial := range input {
		if rawMaterial.Id == uuid.Nil {
			logging.Error("FindRawMaterialById Journey", zap.String("Error", "Invalid uuid"))
			return nil, ErrorIdShouldBeValid
		}

		ids[i] = rawMaterial.Id.String()
	}

	rawMaterialEntities, err := uc.productRepo.FindByIds(ids)
	if err != nil {
		logging.Error("FindRawMaterialById Journey", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error getting raw material: %w", err)
	}

	rawMaterialDTOs := make([]DTOs.RawMaterialDTO, len(rawMaterialEntities))
	for i, rawMaterial := range rawMaterialEntities {
		rawMaterialDTOs[i] = DTOs.MapRawMaterialEntityToDTO(rawMaterial)
	}

	logging.Info("FindRawMaterialById Journey", zap.String("Finish", "FindRawMaterialByIdUseCase"))
	logging.Info("FindRawMaterialById Journey", zap.String("Finish", "FindRawMaterialByIdService"))

	return rawMaterialDTOs, nil
}
