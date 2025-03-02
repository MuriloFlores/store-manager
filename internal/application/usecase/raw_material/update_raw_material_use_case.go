package raw_material

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type updateRawMaterialUseCase struct {
	rawMaterialRepo repositories.RawMaterialsRepositoryInterface
}

type UpdateRawMaterialUseCaseInterface interface {
	UpdateRawMaterial(input []DTOs.RawMaterialDTO) ([]DTOs.RawMaterialDTO, error)
}

func NewUpdateRawMaterialUseCase(rawMaterialRepo repositories.RawMaterialsRepositoryInterface) UpdateRawMaterialUseCaseInterface {
	return &updateRawMaterialUseCase{
		rawMaterialRepo: rawMaterialRepo,
	}
}

func (uc *updateRawMaterialUseCase) UpdateRawMaterial(input []DTOs.RawMaterialDTO) ([]DTOs.RawMaterialDTO, error) {
	logging.Info("UpdateRawMaterial Journey", zap.String("Init", "UpdateRawMaterialUseCase"))
	rawMaterialEntities := make([]entity.RawMaterialInterface, len(input))

	for i, rawMaterial := range input {
		id, err := entity.ParseEntityID(rawMaterial.Id.String())
		if err != nil {
			return nil, ErrorIdShouldBeValid
		}

		rawMaterialEntity := entity.NewRawMaterial(
			&id,
			rawMaterial.Name,
			rawMaterial.Unit,
			rawMaterial.Quantity,
			rawMaterial.Cost.MapMoneyDTOToObject(),
		)

		rawMaterialEntities[i] = rawMaterialEntity
	}

	err := uc.rawMaterialRepo.Update(rawMaterialEntities)
	if err != nil {
		logging.Error("UpdateRawMaterial", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error while updating raw material: %w", err)
	}

	rawMaterialDTOs := make([]DTOs.RawMaterialDTO, len(rawMaterialEntities))
	for i, rawMaterialDTO := range rawMaterialEntities {
		rawMaterialDTOs[i] = DTOs.MapRawMaterialEntityToDTO(rawMaterialDTO)
	}

	logging.Info("UpdateRawMaterial Journey", zap.String("Finish", "UpdateRawMaterialUseCase"))
	logging.Info("UpdateRawMaterial Journey", zap.String("Finish", "UpdateRawMaterialService"))

	return rawMaterialDTOs, nil
}
