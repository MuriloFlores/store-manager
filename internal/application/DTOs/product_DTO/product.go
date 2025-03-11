package product_DTO

import (
	"github.com/google/uuid"
	"store-manager/internal/application/DTOs/money_DTO"
	"store-manager/internal/application/DTOs/raw_material_DTO"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/value_objects"
)

type ProductDTO struct {
	Id             uuid.UUID                         `json:"id"`
	Name           string                            `json:"name"`
	RawMaterials   []raw_material_DTO.RawMaterialDTO `json:"raw_materials"`
	Quantity       int                               `json:"quantity"`
	Value          money_DTO.MoneyDTO                `json:"value"`
	ProductionCost money_DTO.MoneyDTO                `json:"production_cost"`
}

func (p *ProductDTO) MapProductDTOToEntity() entity.ProductInterface {
	id, _ := entity.ParseEntityID(p.Id.String())

	RawMaterialsEntity := make([]entity.RawMaterialInterface, len(p.RawMaterials))

	for i, material := range p.RawMaterials {
		RawMaterialsEntity[i] = material.MapRawMaterialDTOToEntity()
	}

	return entity.NewProduct(
		&id,
		p.Name,
		RawMaterialsEntity,
		p.Quantity,
		p.Value.MapMoneyDTOToObject(),
	)
}

func MapProductEntityToDTO(product entity.ProductInterface) ProductDTO {
	id := uuid.MustParse(product.ID().String())

	rawMaterials := make([]raw_material_DTO.RawMaterialDTO, len(product.RawMaterials()))
	for i, rawMaterial := range product.RawMaterials() {
		rawMaterials[i] = raw_material_DTO.MapRawMaterialEntityToDTO(rawMaterial)
	}

	getMoneyDTO := func(m value_objects.Money) money_DTO.MoneyDTO {
		return money_DTO.MapMoneyObjectToDTO(
			m,
		)
	}

	return ProductDTO{
		Id:             id,
		Name:           product.Name(),
		RawMaterials:   rawMaterials,
		Quantity:       product.Quantity(),
		Value:          getMoneyDTO(product.Value()),
		ProductionCost: getMoneyDTO(product.ProductionCost()),
	}

}
