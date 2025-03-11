package product_assoc_raw_material_DTO

import "github.com/google/uuid"

type ProductAssocRawMaterialDTO struct {
	ProductId    uuid.UUID `json:"product_id"`
	MaterialId   uuid.UUID `json:"material_id"`
	QuantityUsed int       `json:"quantity_used"`
	Activated    bool      `json:"activated"`
}
