package services

import (
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/raw_material_DTO"
	"store-manager/internal/application/usecase/raw_material"
	"store-manager/internal/infrastructure/logging"
)

type rawMaterialService struct {
	createRawMaterialUseCase         raw_material.CreateRawMaterialUseCaseInterface
	findRawMaterialUseCase           raw_material.FindRawMaterialByIdUseCaseInterface
	getAllRawMaterialsUseCase        raw_material.GetAllRawMaterialsUseCaseInterface
	deleteRawMaterialUseCase         raw_material.DeleteRawMaterialUseCaseInterface
	updateRawMaterialUseCase         raw_material.UpdateRawMaterialUseCaseInterface
	getAllRawMaterialsByLimitUseCase raw_material.GetAllRawMaterialsByLimitUseCaseInterface
}

type RawMaterialServiceInterface interface {
	CreateRawMaterial(input []raw_material_DTO.CreateRawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error)
	FindRawMaterialById(input []raw_material_DTO.FindRawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error)
	GetAllRawMaterials() ([]raw_material_DTO.RawMaterialDTO, error)
	GetAllRawMaterialsByLimitRisk() ([]raw_material_DTO.RawMaterialDTO, error)
	DeleteRawMaterial(dtos []raw_material_DTO.FindRawMaterialDTO) error
	UpdateRawMaterial(input []raw_material_DTO.RawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error)
}

func NewRawMaterialService(
	createRawMaterialUseCase raw_material.CreateRawMaterialUseCaseInterface,
	findRawMaterialUseCase raw_material.FindRawMaterialByIdUseCaseInterface,
	getAllRawMaterialsUseCase raw_material.GetAllRawMaterialsUseCaseInterface,
	deleteRawMaterialUseCase raw_material.DeleteRawMaterialUseCaseInterface,
	updateRawMaterialUseCase raw_material.UpdateRawMaterialUseCaseInterface,
	getAllRawMaterialByLimitUseCase raw_material.GetAllRawMaterialsByLimitUseCaseInterface,
) RawMaterialServiceInterface {
	return &rawMaterialService{
		createRawMaterialUseCase:         createRawMaterialUseCase,
		findRawMaterialUseCase:           findRawMaterialUseCase,
		getAllRawMaterialsUseCase:        getAllRawMaterialsUseCase,
		deleteRawMaterialUseCase:         deleteRawMaterialUseCase,
		updateRawMaterialUseCase:         updateRawMaterialUseCase,
		getAllRawMaterialsByLimitUseCase: getAllRawMaterialByLimitUseCase,
	}
}

func (r *rawMaterialService) CreateRawMaterial(input []raw_material_DTO.CreateRawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("CreateRawMaterial Journey", zap.String("Init", "CreateRawMaterialService"))
	return r.createRawMaterialUseCase.CreateRawMaterial(input)
}

func (r *rawMaterialService) FindRawMaterialById(input []raw_material_DTO.FindRawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("FindRawMaterialById Journey", zap.String("Init", "FindRawMaterialByIdService"))
	return r.findRawMaterialUseCase.FindRawMaterialById(input)
}

func (r *rawMaterialService) GetAllRawMaterials() ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("GetAllRawMaterials Journey", zap.String("Init", "GetAllRawMaterialsService"))
	return r.getAllRawMaterialsUseCase.GetAllRawMaterials()
}

func (r *rawMaterialService) GetAllRawMaterialsByLimitRisk() ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("GetAllRawMaterialsByLimit Journey", zap.String("Init", "GetAllRawMaterialsByLimitRisk"))
	return r.getAllRawMaterialsByLimitUseCase.GetAllRawMaterialsByLimit()
}

func (r *rawMaterialService) DeleteRawMaterial(dtos []raw_material_DTO.FindRawMaterialDTO) error {
	logging.Info("DeleteRawMaterial Journey", zap.String("Init", "DeleteRawMaterialService"))
	return r.deleteRawMaterialUseCase.DeleteRawMaterial(dtos)
}

func (r *rawMaterialService) UpdateRawMaterial(input []raw_material_DTO.RawMaterialDTO) ([]raw_material_DTO.RawMaterialDTO, error) {
	logging.Info("UpdateRawMaterial Journey", zap.String("Init", "UpdateRawMaterialService"))
	return r.updateRawMaterialUseCase.UpdateRawMaterial(input)
}
