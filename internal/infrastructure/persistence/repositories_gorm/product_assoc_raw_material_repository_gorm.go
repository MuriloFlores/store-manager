package repositories_gorm

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"store-manager/internal/application/DTOs"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/persistence/models"
)

type ProductAssocRawMaterialRepositoryGorm struct {
	db *gorm.DB
}

var _ repositories.ProductAssocRawMaterialRepositoryInterface

func NewProductAssocRawMaterialRepositoryGorm(db *gorm.DB) repositories.ProductAssocRawMaterialRepositoryInterface {
	return &ProductAssocRawMaterialRepositoryGorm{
		db: db,
	}
}

func (p ProductAssocRawMaterialRepositoryGorm) CreateAssociation(assoc DTOs.ProductAssocRawMaterialDTO) error {
	model := models.ProductRawMaterialModel{
		ProductId:    assoc.ProductId,
		MaterialId:   assoc.MaterialId,
		QuantityUsed: assoc.QuantityUsed,
		Activated:    assoc.Activated,
	}

	return p.db.Create(&model).Error
}

func (p ProductAssocRawMaterialRepositoryGorm) UpdateAssociation(assoc DTOs.ProductAssocRawMaterialDTO) error {
	var model models.ProductRawMaterialModel
	if err := p.db.Where("product_id = ? AND material_id = ?", assoc.ProductId, assoc.MaterialId).First(&model).Error; err != nil {
		return err
	}

	model.QuantityUsed = assoc.QuantityUsed
	model.Activated = assoc.Activated

	return p.db.Save(&model).Error
}

func (p ProductAssocRawMaterialRepositoryGorm) DeleteAssociation(productId, materialId uuid.UUID) error {
	var model models.ProductRawMaterialModel
	if err := p.db.Where("product_id = ? AND material_id = ?", productId, materialId).First(&model).Error; err != nil {
		return err
	}

	return p.db.Delete(&model).Error
}

func (p ProductAssocRawMaterialRepositoryGorm) GetAssociation(productId, materialId uuid.UUID) (DTOs.ProductAssocRawMaterialDTO, error) {
	var model models.ProductRawMaterialModel
	err := p.db.Where("product_id = ? AND material_id = ?", productId, materialId).First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DTOs.ProductAssocRawMaterialDTO{}, nil
		}

		return DTOs.ProductAssocRawMaterialDTO{}, err
	}

	dto := DTOs.ProductAssocRawMaterialDTO{
		ProductId:    model.ProductId,
		MaterialId:   model.MaterialId,
		QuantityUsed: model.QuantityUsed,
		Activated:    model.Activated,
	}

	return dto, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) ListAssociationByProduct(productId uuid.UUID) ([]DTOs.ProductAssocRawMaterialDTO, error) {
	var modelList []models.ProductRawMaterialModel
	err := p.db.Where("product_id = ?", productId).Find(&modelList).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	dtos := make([]DTOs.ProductAssocRawMaterialDTO, len(modelList))
	for i, model := range modelList {
		dtos[i] = DTOs.ProductAssocRawMaterialDTO{
			ProductId:    model.ProductId,
			MaterialId:   model.MaterialId,
			QuantityUsed: model.QuantityUsed,
			Activated:    model.Activated,
		}
	}

	return dtos, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) ListAssociationByRawMaterial(rawMaterial uuid.UUID) ([]DTOs.ProductAssocRawMaterialDTO, error) {
	var modelList []models.ProductRawMaterialModel
	err := p.db.Where("product_id = ?", rawMaterial).Find(&modelList).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	dtos := make([]DTOs.ProductAssocRawMaterialDTO, len(modelList))
	for i, model := range modelList {
		dtos[i] = DTOs.ProductAssocRawMaterialDTO{
			ProductId:    model.ProductId,
			MaterialId:   model.MaterialId,
			QuantityUsed: model.QuantityUsed,
			Activated:    model.Activated,
		}
	}

	return dtos, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) ListAssociationByActivated(activated bool) ([]DTOs.ProductAssocRawMaterialDTO, error) {
	var modelList []models.ProductRawMaterialModel
	err := p.db.Where("activated = ?", activated).Find(&modelList).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	dtos := make([]DTOs.ProductAssocRawMaterialDTO, len(modelList))
	for i, model := range modelList {
		dtos[i] = DTOs.ProductAssocRawMaterialDTO{
			ProductId:    model.ProductId,
			MaterialId:   model.MaterialId,
			QuantityUsed: model.QuantityUsed,
			Activated:    model.Activated,
		}
	}

	return dtos, nil
}
