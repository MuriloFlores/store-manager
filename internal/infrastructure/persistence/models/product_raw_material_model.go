package models

import (
	"github.com/google/uuid"
	"store-manager/internal/application/DTOs/product_assoc_raw_material_DTO"
)

type ProductRawMaterialModel struct {
	ProductId    uuid.UUID `gorm:"type:uuid;"`
	MaterialId   uuid.UUID `gorm:"type:uuid;"`
	QuantityUsed int       `gorm:"type:int;default:0"`
	Activated    bool      `gorm:"type:boolean;default:false"`
}

func (p *ProductRawMaterialModel) MapProductRawMaterialAssocModelToDTO() product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO {
	return product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO{
		ProductId:    p.ProductId,
		MaterialId:   p.MaterialId,
		QuantityUsed: p.QuantityUsed,
		Activated:    p.Activated,
	}
}

func MapProductRawMaterialDTOToModel(assoc product_assoc_raw_material_DTO.ProductAssocRawMaterialDTO) ProductRawMaterialModel {
	return ProductRawMaterialModel{
		ProductId:    assoc.ProductId,
		MaterialId:   assoc.MaterialId,
		QuantityUsed: assoc.QuantityUsed,
		Activated:    assoc.Activated,
	}
}

func (ProductRawMaterialModel) TableName() string {
	return "product_raw_materials"
}
