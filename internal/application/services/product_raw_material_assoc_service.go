package services

import (
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	assoc "store-manager/internal/application/usecase/product_raw_material_assoc"
	"store-manager/internal/domain/entity"
	"store-manager/internal/infrastructure/logging"
)

type productRawMaterialAssocService struct {
	createAssoc         assoc.CreateAssocUseCaseInterface
	deleteAssoc         assoc.DeleteAssocUseCaseInterface
	findAllAssoc        assoc.FindAllAssociationsUseCaseInterface
	findAssocByCriteria assoc.FindAssocByCriteriaUseCaseInterface
	findAssocById       assoc.FindAssocByIdUseCaseInterface
	updateAssoc         assoc.UpdateAssocUseCaseInterface
}

type ProductRawMaterialAssocServiceInterface interface {
	CreateAssoc(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error)
	DeleteAssoc(productIds []string, materialsIds []string) error
	FindAllAssociations() ([]entity.ProductInterface, error)
	FindByCriteria(criteria string, info interface{}) ([]entity.ProductInterface, error)
	FindAssocById(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error)
	UpdateAssoc(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error)
}

func NewProductRawMaterialAssocService(
	createAssoc assoc.CreateAssocUseCaseInterface,
	deleteAssoc assoc.DeleteAssocUseCaseInterface,
	findAllAssoc assoc.FindAllAssociationsUseCaseInterface,
	findAssocByCriteria assoc.FindAssocByCriteriaUseCaseInterface,
	findAssocById assoc.FindAssocByIdUseCaseInterface,
	updateAssoc assoc.UpdateAssocUseCaseInterface,
) ProductRawMaterialAssocServiceInterface {
	return &productRawMaterialAssocService{
		createAssoc:         createAssoc,
		deleteAssoc:         deleteAssoc,
		findAllAssoc:        findAllAssoc,
		findAssocByCriteria: findAssocByCriteria,
		findAssocById:       findAssocById,
		updateAssoc:         updateAssoc,
	}
}

func (p *productRawMaterialAssocService) CreateAssoc(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error) {
	logging.Info("CreateAssoc Journey", zap.String("Init", "CreateAssocService"))
	return p.createAssoc.CreateAssoc(input)
}

func (p *productRawMaterialAssocService) DeleteAssoc(productIds []string, materialsIds []string) error {
	logging.Info("DeleteAssoc Journey", zap.String("Init", "DeleteAssocService"))
	return p.deleteAssoc.DeleteByIds(productIds, materialsIds)
}

func (p *productRawMaterialAssocService) FindAllAssociations() ([]entity.ProductInterface, error) {
	logging.Info("FindAllAssoc Journey", zap.String("Init", "FindAllAssocService"))
	return p.findAllAssoc.GetAllAssociations()
}

func (p *productRawMaterialAssocService) FindByCriteria(criteria string, info interface{}) ([]entity.ProductInterface, error) {
	logging.Info("FindByCriteria Journey", zap.String("Init", "FindByCriteriaService"))
	return p.findAssocByCriteria.FindByCriteria(criteria, info)
}

func (p *productRawMaterialAssocService) FindAssocById(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error) {
	logging.Info("FindAssocById Journey", zap.String("Init", "FindAssocByIdService"))
	return p.findAssocById.FindByIds(input)
}

func (p *productRawMaterialAssocService) UpdateAssoc(input []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ([]entity.ProductInterface, error) {
	logging.Info("UpdateAssoc Journey", zap.String("Init", "UpdateAssocService"))
	return p.updateAssoc.UpdateByIds(input)
}
