package repositories

import (
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
	"store-manager/internal/domain/entity"
)

type ProductAssocRawMaterialRepositoryInterface interface {
	CreateAssociation(assoc []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) error
	UpdateAssociation(assoc []product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) error
	DeleteAssociation(productId, materialId []string) error
	GetAssociation(productId, materialId []string) ([]entity.ProductInterface, error)
	GetAllAssociations() ([]entity.ProductInterface, error)
	ListAssociationByProduct(productId []string) ([]entity.ProductInterface, error)
	ListAssociationByRawMaterial(rawMaterialId []string) ([]entity.ProductInterface, error)
	ListAssociationByActivated(activated bool) ([]entity.ProductInterface, error)
}
