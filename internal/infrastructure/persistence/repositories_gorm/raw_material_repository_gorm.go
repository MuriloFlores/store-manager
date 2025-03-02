package repositories_gorm

import (
	"errors"
	"gorm.io/gorm"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/persistence/models"
)

type rawMaterialsRepositoryGorm struct {
	db *gorm.DB
}

var _ repositories.RawMaterialsRepositoryInterface = (*rawMaterialsRepositoryGorm)(nil)

func NewRawMaterialsRepositoryGorm(db *gorm.DB) repositories.RawMaterialsRepositoryInterface {
	return &rawMaterialsRepositoryGorm{db: db}
}

func (r *rawMaterialsRepositoryGorm) Save(rawMaterials []entity.RawMaterialInterface) error {
	if len(rawMaterials) == 0 {
		return errors.New("rawMaterials is empty")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		rawMaterialsModels := make([]models.RawMaterialModel, len(rawMaterials))
		for i, rawMaterial := range rawMaterials {
			rawMaterialsModels[i] = models.MapRawMaterialEntityToModel(rawMaterial)
		}

		if err := tx.Create(&rawMaterialsModels).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *rawMaterialsRepositoryGorm) FindByIds(ids []string) ([]entity.RawMaterialInterface, error) {
	var rawMaterialModels []models.RawMaterialModel

	err := r.db.Where("id IN (?)", ids).Find(&rawMaterialModels).Error
	if err != nil {
		return nil, err
	}

	if len(rawMaterialModels) == 0 {
		return []entity.RawMaterialInterface{}, ErrorRecordNotFound
	}

	rawMaterialEntities := make([]entity.RawMaterialInterface, len(rawMaterialModels))
	for i, model := range rawMaterialModels {
		rawMaterialEntities[i] = model.MapRawMaterialModelToEntity()
	}

	return rawMaterialEntities, nil
}

func (r *rawMaterialsRepositoryGorm) Update(rawMaterials []entity.RawMaterialInterface) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, rawMaterial := range rawMaterials {
			model := models.MapRawMaterialEntityToModel(rawMaterial)
			if err := tx.Save(model).Updates(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *rawMaterialsRepositoryGorm) DeleteByIds(ids []string) error {
	result := r.db.Delete(&models.RawMaterialModel{}, "id IN (?)", ids)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrorRecordNotFound
	}
	return nil
}

func (r *rawMaterialsRepositoryGorm) GetAllRawMaterials() ([]entity.RawMaterialInterface, error) {
	var rawMaterialModels []models.RawMaterialModel

	err := r.db.
		Order("quantity ASC").
		Find(&rawMaterialModels).Error

	if err != nil {
		return nil, err
	}

	if len(rawMaterialModels) == 0 {
		return nil, ErrorRecordNotFound
	}

	rawMaterialsEntities := make([]entity.RawMaterialInterface, len(rawMaterialModels))
	for i, rawMaterial := range rawMaterialModels {
		rawMaterialsEntities[i] = rawMaterial.MapRawMaterialModelToEntity()
	}

	return rawMaterialsEntities, err
}

func (r *rawMaterialsRepositoryGorm) GetAllRawMaterialsByLimitRisk() ([]entity.RawMaterialInterface, error) {
	var rawMaterialModels []models.RawMaterialModel

	err := r.db.
		Order("ABS(risk_limit - quantity) ASC").
		Find(&rawMaterialModels).Error

	if err != nil {
		return nil, err
	}

	if len(rawMaterialModels) == 0 {
		return nil, ErrorRecordNotFound
	}

	rawMaterialsEntities := make([]entity.RawMaterialInterface, len(rawMaterialModels))
	for i, rawMaterial := range rawMaterialModels {
		rawMaterialsEntities[i] = rawMaterial.MapRawMaterialModelToEntity()
	}

	return rawMaterialsEntities, err
}
