package models

import (
	"github.com/google/uuid"
	"store-manager/internal/domain/entity"
)

type RawMaterialModel struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name      string     `gorm:"type:varchar(255);not null"`
	Unit      string     `gorm:"type:varchar(20);not null"`
	Quantity  int        `gorm:"default:0"`
	Cost      MoneyModel `gorm:"embedded;embeddedPrefix:cost_"`
	RiskLimit int        `gorm:"default:0"`
}

func (r *RawMaterialModel) MapRawMaterialModelToEntity() entity.RawMaterialInterface {
	id, _ := entity.ParseEntityID(r.Id.String())

	return entity.NewRawMaterial(
		&id,
		r.Name,
		entity.Unit(r.Unit),
		r.Quantity,
		r.Cost.MapMoneyModelToEntity(),
		&r.RiskLimit,
	)
}

func MapRawMaterialEntityToModel(rawMaterial entity.RawMaterialInterface) RawMaterialModel {
	id := uuid.MustParse(rawMaterial.ID().String())
	cost := MapMoneyObjectToMoneyModel(rawMaterial.Cost())

	return RawMaterialModel{
		Id:       id,
		Name:     rawMaterial.Name(),
		Unit:     string(rawMaterial.Unit()),
		Quantity: rawMaterial.Quantity(),
		Cost:     cost,
	}
}
