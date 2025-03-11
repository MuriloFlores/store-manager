package product_assoc_raw_material_DTO

import "github.com/google/uuid"

type FindProductAssocRawMaterialDTO struct {
	ProductIds  []uuid.UUID `json:"product_ids"`
	MaterialIds []uuid.UUID `json:"material_ids"`
}
