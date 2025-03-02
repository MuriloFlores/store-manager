package services

import (
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/application/usecase/raw_material"
	"store-manager/internal/infrastructure/logging"
)

type rawMaterialService struct {
	createRawMaterialUseCase  raw_material.CreateRawMaterialUseCaseInterface
	findRawMaterialUseCase    raw_material.FindRawMaterialByIdUseCaseInterface
	getAllRawMaterialsUseCase raw_material.GetAllRawMaterialsUseCaseInterface
	deleteRawMaterialUseCase  raw_material.DeleteRawMaterialUseCaseInterface
	updateRawMaterialUseCase  raw_material.UpdateRawMaterialUseCaseInterface
}

type RawMaterialServiceInterface interface {
	CreateRawMaterial(input []DTOs.CreateRawMaterialDTO) ([]DTOs.RawMaterialDTO, error)
	FindRawMaterialById(input []DTOs.FindRawMaterialDTO) ([]DTOs.RawMaterialDTO, error)
	GetAllRawMaterials() ([]DTOs.RawMaterialDTO, error)
	DeleteRawMaterial(dtos []DTOs.FindRawMaterialDTO) error
	UpdateRawMaterial(input []DTOs.RawMaterialDTO) ([]DTOs.RawMaterialDTO, error)
}

func NewRawMaterialService(
	createRawMaterialUseCase raw_material.CreateRawMaterialUseCaseInterface,
	findRawMaterialUseCase raw_material.FindRawMaterialByIdUseCaseInterface,
	getAllRawMaterialsUseCase raw_material.GetAllRawMaterialsUseCaseInterface,
	deleteRawMaterialUseCase raw_material.DeleteRawMaterialUseCaseInterface,
	updateRawMaterialUseCase raw_material.UpdateRawMaterialUseCaseInterface,
) RawMaterialServiceInterface {
	return &rawMaterialService{
		createRawMaterialUseCase:  createRawMaterialUseCase,
		findRawMaterialUseCase:    findRawMaterialUseCase,
		getAllRawMaterialsUseCase: getAllRawMaterialsUseCase,
		deleteRawMaterialUseCase:  deleteRawMaterialUseCase,
		updateRawMaterialUseCase:  updateRawMaterialUseCase,
	}
}

func (r *rawMaterialService) CreateRawMaterial(input []DTOs.CreateRawMaterialDTO) ([]DTOs.RawMaterialDTO, error) {
	logging.Info("CreateRawMaterial Journey", zap.String("Init", "CreateRawMaterialService"))
	return r.createRawMaterialUseCase.CreateRawMaterial(input)
}

func (r *rawMaterialService) FindRawMaterialById(input []DTOs.FindRawMaterialDTO) ([]DTOs.RawMaterialDTO, error) {
	logging.Info("FindRawMaterialById Journey", zap.String("Init", "FindRawMaterialByIdService"))
	return r.findRawMaterialUseCase.FindRawMaterialById(input)
}

func (r *rawMaterialService) GetAllRawMaterials() ([]DTOs.RawMaterialDTO, error) {
	logging.Info("GetAllRawMaterials Journey", zap.String("Init", "GetAllRawMaterialsService"))
	return r.getAllRawMaterialsUseCase.GetAllRawMaterials()
}

func (r *rawMaterialService) DeleteRawMaterial(dtos []DTOs.FindRawMaterialDTO) error {
	logging.Info("DeleteRawMaterial Journey", zap.String("Init", "DeleteRawMaterialService"))
	return r.deleteRawMaterialUseCase.DeleteRawMaterial(dtos)
}

func (r *rawMaterialService) UpdateRawMaterial(input []DTOs.RawMaterialDTO) ([]DTOs.RawMaterialDTO, error) {
	logging.Info("UpdateRawMaterial Journey", zap.String("Init", "UpdateRawMaterialService"))
	return r.updateRawMaterialUseCase.UpdateRawMaterial(input)
}
