package repositories_gorm

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/persistence/models"
)

type ProductAssocRawMaterialRepositoryGorm struct {
	db          *gorm.DB
	productRepo repositories.ProductRepositoryInterface
	rawRepo     repositories.RawMaterialsRepositoryInterface
}

var _ repositories.ProductAssocRawMaterialRepositoryInterface

func NewProductAssocRawMaterialRepositoryGorm(db *gorm.DB, productRepo repositories.ProductRepositoryInterface, rawRepo repositories.RawMaterialsRepositoryInterface) repositories.ProductAssocRawMaterialRepositoryInterface {
	return &ProductAssocRawMaterialRepositoryGorm{
		db:          db,
		productRepo: productRepo,
		rawRepo:     rawRepo,
	}
}

func (p ProductAssocRawMaterialRepositoryGorm) CreateAssociation(associations []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) error {
	if len(associations) == 0 {
		return nil
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		assocModels := make([]models.ProductRawMaterialModel, len(associations))
		for i, assoc := range associations {
			assocModels[i] = models.MapProductRawMaterialDTOToModel(assoc)
		}

		if err := tx.Create(&assocModels).Error; err != nil {
			return fmt.Errorf("error creating associations: %w", err)
		}

		productMaterialsMap := make(map[uuid.UUID]map[uuid.UUID]struct{})
		for _, assocModel := range assocModels {
			prodID := assocModel.ProductId
			matID := assocModel.MaterialId

			if productMaterialsMap[prodID] == nil {
				productMaterialsMap[prodID] = make(map[uuid.UUID]struct{})
			}

			productMaterialsMap[prodID][matID] = struct{}{}
		}

		var uniqueProductIDs []string
		for prodID := range productMaterialsMap {
			uniqueProductIDs = append(uniqueProductIDs, prodID.String())
		}

		products, err := p.productRepo.FindByIds(uniqueProductIDs)
		if err != nil {
			return fmt.Errorf("error fetching products: %w", err)
		}
		if len(products) == 0 {
			return nil
		}

		for _, prod := range products {
			prodID := prod.ID()

			if materialSet, exists := productMaterialsMap[uuid.MustParse(prodID.String())]; exists {

				var rawMaterialIDs []string
				for matID := range materialSet {
					rawMaterialIDs = append(rawMaterialIDs, matID.String())
				}

				rawMaterials, err := p.rawRepo.FindByIds(rawMaterialIDs)
				if err != nil {
					return fmt.Errorf("error fetching raw materials for product %s: %w", prodID.String(), err)
				}

				prod.AddRawMaterials(rawMaterials)

				prod.CalculateCost()

				updatedCost := prod.ProductionCost()

				fmt.Println(prod.ProductionCost())
				fmt.Println(updatedCost.ValueInCents())

				if err := tx.Model(&models.ProductModel{}).
					Where("id = ?", prodID.String()).
					Updates(map[string]interface{}{
						"production_cost_total_in_cents": updatedCost.ValueInCents(),
						"production_cost_currency":       updatedCost.Currency(),
					}).Error; err != nil {
					return fmt.Errorf("error updating production cost for product %s: %w", prodID.String(), err)
				}
			}
		}

		return nil
	})
}

func (p ProductAssocRawMaterialRepositoryGorm) UpdateAssociation(associations []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, assoc := range associations {
			model := models.MapProductRawMaterialDTOToModel(assoc)
			if err := tx.Save(&model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p ProductAssocRawMaterialRepositoryGorm) DeleteAssociation(productIds, materialIds []string) error {
	fmt.Println("productIds: ", productIds)
	fmt.Println("materialIds: ", materialIds)

	return p.db.
		Where("product_id IN (?) AND material_id IN (?)", productIds, materialIds).
		Delete(&models.ProductRawMaterialModel{}).Error
}

func (p ProductAssocRawMaterialRepositoryGorm) GetAssociation(productIds, materialIds []string) ([]entity.ProductInterface, error) {
	var assocModels []models.ProductRawMaterialModel
	err := p.db.
		Where("product_id IN (?) AND material_id IN (?)", productIds, materialIds).
		Find(&assocModels).Error
	if err != nil {
		return nil, err
	}
	if len(assocModels) == 0 {
		return []entity.ProductInterface{}, nil
	}

	groups := make(map[string]map[uuid.UUID]struct{})
	for _, assoc := range assocModels {
		prodID := assoc.ProductId.String()
		if groups[prodID] == nil {
			groups[prodID] = make(map[uuid.UUID]struct{})
		}
		groups[prodID][assoc.MaterialId] = struct{}{}
	}

	var uniqueProductIDs []string
	for prodID := range groups {
		uniqueProductIDs = append(uniqueProductIDs, prodID)
	}

	products, err := p.productRepo.FindByIds(uniqueProductIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}

	for _, prod := range products {
		prodID := prod.ID().String()
		if materialSet, exists := groups[prodID]; exists {
			var rawMaterialIDs []string
			for matID := range materialSet {
				rawMaterialIDs = append(rawMaterialIDs, matID.String())
			}
			rawMaterials, err := p.rawRepo.FindByIds(rawMaterialIDs)
			if err != nil {
				return nil, fmt.Errorf("error fetching raw materials for product %s: %w", prodID, err)
			}

			prod.AddRawMaterials(rawMaterials)

			prod.CalculateCost()
		}
	}

	return products, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) ListAssociationByProduct(productId []string) ([]entity.ProductInterface, error) {
	var assocModels []models.ProductRawMaterialModel

	err := p.db.Where("product_id IN (?)", productId).Find(&assocModels).Error
	if err != nil {
		return nil, err
	}

	if len(assocModels) == 0 {
		return []entity.ProductInterface{}, nil
	}

	groups := make(map[string]map[string]struct{})
	for _, assoc := range assocModels {
		prodId := assoc.ProductId.String()

		if groups[prodId] == nil {
			groups[prodId] = make(map[string]struct{})
		}

		groups[prodId][assoc.MaterialId.String()] = struct{}{}
	}

	var uniqueProductsIds []string
	for prodId := range groups {
		uniqueProductsIds = append(uniqueProductsIds, prodId)
	}

	products, err := p.productRepo.FindByIds(uniqueProductsIds)
	if err != nil {
		return []entity.ProductInterface{}, fmt.Errorf("error fetching products: %w", err)
	}

	for _, prod := range products {
		prodId := prod.ID().String()

		if materialSet, exists := groups[prodId]; exists {
			var rawMatIDs []string

			for matID := range materialSet {
				rawMatIDs = append(rawMatIDs, matID)
			}

			rawMaterials, err := p.rawRepo.FindByIds(rawMatIDs)
			if err != nil {
				return []entity.ProductInterface{}, fmt.Errorf("error fetching materials: %w", err)
			}

			prod.AddRawMaterials(rawMaterials)
		}
	}

	return products, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) ListAssociationByRawMaterial(rawMaterial []string) ([]entity.ProductInterface, error) {
	var assocModels []models.ProductRawMaterialModel
	err := p.db.Where("material_id IN (?)", rawMaterial).Find(&assocModels).Error
	if err != nil {
		return nil, err
	}

	if len(assocModels) == 0 {
		return []entity.ProductInterface{}, nil
	}

	groups := make(map[string]map[string]struct{})
	for _, assoc := range assocModels {
		prodId := assoc.ProductId.String()

		if groups[prodId] == nil {
			groups[prodId] = make(map[string]struct{})
		}

		groups[prodId][assoc.MaterialId.String()] = struct{}{}
	}

	var uniqueProductsIds []string
	for prodId := range groups {
		uniqueProductsIds = append(uniqueProductsIds, prodId)
	}

	products, err := p.productRepo.FindByIds(uniqueProductsIds)
	if err != nil {
		return []entity.ProductInterface{}, fmt.Errorf("error fetching products: %w", err)
	}

	for _, prod := range products {
		prodId := prod.ID().String()

		if materialSet, exists := groups[prodId]; exists {
			var rawMatIDs []string

			for matID := range materialSet {
				rawMatIDs = append(rawMatIDs, matID)
			}

			rawMaterials, err := p.rawRepo.FindByIds(rawMatIDs)
			if err != nil {
				return []entity.ProductInterface{}, fmt.Errorf("error fetching materials: %w", err)
			}

			prod.AddRawMaterials(rawMaterials)
		}
	}

	return products, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) ListAssociationByActivated(activated bool) ([]entity.ProductInterface, error) {
	var assocModels []models.ProductRawMaterialModel

	err := p.db.Where("activated = ?", activated).Find(&assocModels).Error
	if err != nil {
		return nil, err
	}

	if len(assocModels) == 0 {
		return []entity.ProductInterface{}, nil
	}

	groups := make(map[string]map[string]struct{})
	for _, assoc := range assocModels {
		prodId := assoc.ProductId.String()

		if groups[prodId] == nil {
			groups[prodId] = make(map[string]struct{})
		}

		groups[prodId][assoc.MaterialId.String()] = struct{}{}
	}

	var uniqueProductsIds []string
	for prodId := range groups {
		uniqueProductsIds = append(uniqueProductsIds, prodId)
	}

	products, err := p.productRepo.FindByIds(uniqueProductsIds)
	if err != nil {
		return []entity.ProductInterface{}, fmt.Errorf("error fetching products: %w", err)
	}

	for _, prod := range products {
		prodId := prod.ID().String()

		if materialSet, exists := groups[prodId]; exists {
			var rawMatIDs []string

			for matID := range materialSet {
				rawMatIDs = append(rawMatIDs, matID)
			}

			rawMaterials, err := p.rawRepo.FindByIds(rawMatIDs)
			if err != nil {
				return []entity.ProductInterface{}, fmt.Errorf("error fetching materials: %w", err)
			}

			prod.AddRawMaterials(rawMaterials)
		}
	}

	return products, nil
}

func (p ProductAssocRawMaterialRepositoryGorm) GetAllAssociations() ([]entity.ProductInterface, error) {
	var assocModels []models.ProductRawMaterialModel

	err := p.db.Find(&assocModels).Error
	if err != nil {
		return []entity.ProductInterface{}, err
	}

	if len(assocModels) == 0 {
		return []entity.ProductInterface{}, nil
	}

	groups := make(map[string]map[string]struct{})
	for _, assoc := range assocModels {
		prodId := assoc.ProductId.String()

		if groups[prodId] == nil {
			groups[prodId] = make(map[string]struct{})
		}

		groups[prodId][assoc.MaterialId.String()] = struct{}{}
	}

	var uniqueProductsIds []string
	for prodId := range groups {
		uniqueProductsIds = append(uniqueProductsIds, prodId)
	}

	products, err := p.productRepo.FindByIds(uniqueProductsIds)
	if err != nil {
		return []entity.ProductInterface{}, fmt.Errorf("error fetching products: %w", err)
	}

	for _, prod := range products {
		prodId := prod.ID().String()
		if materialSet, exists := groups[prodId]; exists {
			var rawMatIDs []string
			for matID := range materialSet {
				rawMatIDs = append(rawMatIDs, matID)
			}

			rawMaterials, err := p.rawRepo.FindByIds(rawMatIDs)
			if err != nil {
				return []entity.ProductInterface{}, fmt.Errorf("error fetching materials: %w", err)
			}

			prod.AddRawMaterials(rawMaterials)
		}
	}

	return products, nil
}
