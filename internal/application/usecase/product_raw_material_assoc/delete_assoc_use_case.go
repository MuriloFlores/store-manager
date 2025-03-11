package product_raw_material_assoc

import "store-manager/internal/domain/repositories"

type deleteIdProductRawMaterialAssocUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

func NewDeleteByIdProduct