package raw_material

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/raw_material_DTO"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type getAllRawMaterialsByLimitUseCase struct {
	rawMaterialRepo repositories.RawMaterialsRepositoryInterface
}

type GetAllRawMaterialsByLimitUseCaseInterface interface {
	GetAllRawMaterialsByLimit() ([]raw_material_DTO.RawMaterialDTO, error)
}

func NewGetAllRawMaterialsByLimit(rawMaterialRepo repositories.RawMaterialsRepositoryInterface) GetAllRawMaterialsByLimitUseCaseInterface {
	return &getAllRawMaterialsByLimitUseCase{
		rawMaterialRepo: rawMaterialRepo,
	}
}

func (u *getAllRawMaterialsByLimitUseCase) GetAllRawMaterialsByLimit() ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("GetAllRawMaterialsByLimit Journey", zap.String("Init", "GetAllRawMaterialsByLimitUseCase"))

	rawMaterialsEntities, err := u.rawMaterialRepo.GetAllRawMaterialsByLimitRisk()
	if err != nil {
		logging.Error("GetAllRawMaterials Journey", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error while getting all raw materials: %w", err)
	}

	rawMaterialDTOs := make([]raw_material_DTO.RawMaterialDTO, len(rawMaterialsEntities))
	for i, rawMaterial := range rawMaterialsEntities {
		rawMaterialDTOs[i] = raw_material_DTO.MapRawMaterialEntityToDTO(rawMaterial)
	}

	logging.Info("GetAllRawMaterialsByLimit Journey", zap.String("Finish", "GetAllRawMaterialsByLimitUseCase"))
	logging.Info("GetAllRawMaterialsByLimit Journey", zap.String("Finish", "GetAllRawMaterialsByLimitService"))

	return rawMaterialDTOs, nil
}
