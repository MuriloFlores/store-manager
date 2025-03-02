package models

import (
	"github.com/google/uuid"
	"store-manager/internal/domain/entity"
)

type ProductModel struct {
	Id             uuid.UUID  `gorm:"type:uuid;primaryKey;"`
	Name           string     `gorm:"type:varchar(255);not null"`
	Quantity       int        `gorm:"type:int;not null"`
	Value          MoneyModel `gorm:"embedded;embeddedPrefix:value_"`
	ProductionCost MoneyModel `gorm:"embedded;embeddedPrefix:production_cost_"`
}

func (p *ProductModel) MapProductModelToEntity() entity.ProductInterface {
	id, _ := entity.ParseEntityID(p.Id.String())

	return entity.NewProduct(
		&id,
		p.Name,
		[]entity.RawMaterialInterface{},
		p.Quantity,
		p.Value.MapMoneyModelToEntity(),
	)
}

func MapProductEntityToModel(product entity.ProductInterface) ProductModel {
	rawMaterialsModels := make([]RawMaterialModel, len(product.RawMaterials()))

	for i, rawMaterial := range product.RawMaterials() {
		rawMaterialsModels[i] = MapRawMaterialEntityToModel(rawMaterial)
	}

	return ProductModel{
		Id:             uuid.MustParse(product.ID().String()),
		Name:           product.Name(),
		Quantity:       product.Quantity(),
		Value:          MapMoneyObjectToMoneyModel(product.Value()),
		ProductionCost: MapMoneyObjectToMoneyModel(product.ProductionCost()),
	}
}
