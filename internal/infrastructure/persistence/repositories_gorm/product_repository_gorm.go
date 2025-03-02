package repositories_gorm

import (
	"errors"
	"gorm.io/gorm"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/persistence/models"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
)

type ProductRepositoryGorm struct {
	db *gorm.DB
}

var _ repositories.ProductRepositoryInterface = (*ProductRepositoryGorm)(nil)

func NewProductRepositoryGorm(db *gorm.DB) repositories.ProductRepositoryInterface {
	return &ProductRepositoryGorm{db: db}
}

func (r *ProductRepositoryGorm) Save(products []entity.ProductInterface) error {
	if len(products) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		productModels := make([]models.ProductModel, len(products))
		for i, product := range products {
			productModels[i] = models.MapProductEntityToModel(product)
		}

		if err := tx.Create(&productModels).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProductRepositoryGorm) GetAllProducts() ([]entity.ProductInterface, error) {
	var productModels []models.ProductModel

	err := r.db.Find(&productModels).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}

		return nil, err
	}

	productEntities := make([]entity.ProductInterface, len(productModels))
	for i, productModel := range productModels {
		productEntities[i] = productModel.MapProductModelToEntity()
	}

	return productEntities, nil
}

func (r *ProductRepositoryGorm) FindByIds(ids []string) ([]entity.ProductInterface, error) {
	var productModels []models.ProductModel

	err := r.db.Find(&productModels, "id IN ?", ids).Error
	if err != nil {
		return nil, err
	}

	if len(productModels) == 0 {
		return []entity.ProductInterface{}, ErrorRecordNotFound
	}

	productEntities := make([]entity.ProductInterface, len(productModels))
	for i, model := range productModels {
		productEntities[i] = model.MapProductModelToEntity()
	}

	return productEntities, nil
}

func (r *ProductRepositoryGorm) Update(products []entity.ProductInterface) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			model := models.MapProductEntityToModel(product)
			if err := tx.Save(&model).Updates(model).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ProductRepositoryGorm) DeleteByIds(ids []string) error {
	err := r.db.Delete(&models.ProductModel{}, "id IN ?", ids).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorRecordNotFound
		}
		return err
	}

	return nil
}
