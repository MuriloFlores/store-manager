package repositories

import (
	"github.com/google/uuid"
	"store-manager/internal/application/DTOs"
)

type ProductAssocRawMaterialRepositoryInterface interface {
	CreateAssociation(assoc DTOs.ProductAssocRawMaterialDTO) error
	UpdateAssociation(assoc DTOs.ProductAssocRawMaterialDTO) error
	DeleteAssociation(productId, materialId uuid.UUID) error
	GetAssociation(productId, materialId uuid.UUID) (DTOs.ProductAssocRawMaterialDTO, error)
	ListAssociationByProduct(productId uuid.UUID) ([]DTOs.ProductAssocRawMaterialDTO, error)
	ListAssociationByRawMaterial(rawMaterial uuid.UUID) ([]DTOs.ProductAssocRawMaterialDTO, error)
	ListAssociationByActivated(activated bool) ([]DTOs.ProductAssocRawMaterialDTO, error)
}
