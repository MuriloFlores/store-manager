package raw_material

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/raw_material_DTO"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type getAllRawMaterialsUseCase struct {
	rawMaterialRepo repositories.RawMaterialsRepositoryInterface
}

type GetAllRawMaterialsUseCaseInterface interface {
	GetAllRawMaterials() ([]raw_material_DTO.RawMaterialDTO, error)
}

func NewGetAllRawMaterials(rawMaterialRepo repositories.RawMaterialsRepositoryInterface) GetAllRawMaterialsUseCaseInterface {
	return &getAllRawMaterialsUseCase{
		rawMaterialRepo: rawMaterialRepo,
	}
}

func (u *getAllRawMaterialsUseCase) GetAllRawMaterials() ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("GetAllRawMaterials Journey", zap.String("Init", "GetAllRawMaterialsUseCase"))

	rawMaterialsEntities, err := u.rawMaterialRepo.GetAllRawMaterials()
	if err != nil {
		logging.Error("GetAllRawMaterials Journey", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("error while getting all raw materials: %w", err)
	}

	rawMaterialDTOs := make([]raw_material_DTO.RawMaterialDTO, len(rawMaterialsEntities))
	for i, rawMaterial := range rawMaterialsEntities {
		rawMaterialDTOs[i] = raw_material_DTO.MapRawMaterialEntityToDTO(rawMaterial)
	}

	logging.Info("GetAllRawMaterials Journey", zap.String("Finish", "GetAllRawMaterialsUseCase"))
	logging.Info("GetAllRawMaterials Journey", zap.String("Finish", "GetAllRawMaterialsService"))

	return rawMaterialDTOs, nil
}
