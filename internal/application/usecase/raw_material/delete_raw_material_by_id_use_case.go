package raw_material

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/raw_material_DTO"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type deleteRawMaterialUseCase struct {
	rawMaterialRepo repositories.RawMaterialsRepositoryInterface
}

var (
	ErrorIdShouldBeValid = errors.New("id should be valid")
)

type DeleteRawMaterialUseCaseInterface interface {
	DeleteRawMaterial(dtos []raw_material_DTO.FindRawMaterialDTO) error
}

func NewDeleteRawMaterialUseCase(rawMaterialRepo repositories.RawMaterialsRepositoryInterface) DeleteRawMaterialUseCaseInterface {
	return &deleteRawMaterialUseCase{
		rawMaterialRepo: rawMaterialRepo,
	}
}

func (uc *deleteRawMaterialUseCase) DeleteRawMaterial(dtos []raw_material_DTO.FindRawMaterialDTO) error {
	logging.Info("DeleteRawMaterial Journey", zap.String("Init", "DeleteRawMaterialUseCase"))

	ids := make([]string, len(dtos))
	for i, rawMaterial := range dtos {
		if rawMaterial.Id == uuid.Nil {
			logging.Error("DeleteRawMaterial Journey", zap.String("Error", "Invalid uuid"))
			return ErrorIdShouldBeValid
		}

		ids[i] = rawMaterial.Id.String()
	}

	err := uc.rawMaterialRepo.DeleteByIds(ids)
	if err != nil {
		logging.Error("DeleteRawMaterial Journey", zap.String("Error", err.Error()))
		return fmt.Errorf("error while deleting raw materials: %w", err)
	}

	logging.Info("DeleteRawMaterial Journey", zap.String("Finish", "DeleteRawMaterialUseCase"))
	logging.Info("DeleteRawMaterial Journey", zap.String("Finish", "DeleteRawMaterialService"))

	return nil
}
