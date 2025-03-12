package raw_material

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/raw_material_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type createRawMaterialUseCase struct {
	rawMaterialRepo repositories.RawMaterialsRepositoryInterface
}

type CreateRawMaterialUseCaseInterface interface {
	CreateRawMaterial(input []raw_material_DTO.CreateRawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error)
}

func NewCreateRawMaterialUseCase(rawMaterialRepo repositories.RawMaterialsRepositoryInterface) CreateRawMaterialUseCaseInterface {
	return &createRawMaterialUseCase{
		rawMaterialRepo: rawMaterialRepo,
	}
}

func (uc *createRawMaterialUseCase) CreateRawMaterial(input []raw_material_DTO.CreateRawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("CreateRawMaterial Journey", zap.String("Init", "CreateRawMaterialUseCase"))
	rawMaterialEntities := make([]entity.RawMaterialInterface, len(input))

	for i, rawMaterial := range input {

		rawMaterialEntity := entity.NewRawMaterial(
			nil,
			rawMaterial.Name,
			rawMaterial.Unit,
			rawMaterial.Quantity,
			rawMaterial.Cost.MapMoneyDTOToObject(),
			nil,
		)

		rawMaterialEntities[i] = rawMaterialEntity
	}

	err := uc.rawMaterialRepo.Save(rawMaterialEntities)
	if err != nil {
		logging.Error("CreateRawMaterial Journey", zap.String("error", err.Error()))
		return nil, fmt.Errorf("error saving raw material %w", err)
	}

	rawMaterialDTOs := make([]raw_material_DTO.RawMaterialDTO, len(input))

	for i, rawMaterial := range rawMaterialEntities {
		rawMaterialDTOs[i] = raw_material_DTO.MapRawMaterialEntityToDTO(rawMaterial)
	}

	logging.Info("CreateRawMaterial Journey", zap.String("Finish", "CreateRawMaterialUseCase"))
	logging.Info("CreateRawMaterial Journey", zap.String("Finish", "CreateRawMaterialService"))

	return rawMaterialDTOs, nil
}
