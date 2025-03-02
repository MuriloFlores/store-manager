package models

import (
	"github.com/google/uuid"
)

type ProductRawMaterialModel struct {
	ProductId    uuid.UUID `gorm:"type:uuid;"`
	MaterialId   uuid.UUID `gorm:"type:uuid;"`
	QuantityUsed int       `gorm:"type:int;default:0"`
	Activated    bool      `gorm:"type:boolean;default:false"`
}

func (ProductRawMaterialModel) TableName() string {
	return "product_raw_materials"
}
