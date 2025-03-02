package DTOs

import (
	"github.com/google/uuid"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/value_objects"
)

type UpdateProductDTO struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Value    MoneyDTO  `json:"value"`
}

type CreateProductDTO struct {
	Name     string   `json:"name"`
	Quantity int      `json:"quantity"`
	Value    MoneyDTO `json:"value"`
}

type FindProductDTO struct {
	Id uuid.UUID `json:"id"`
}

type ProductDTO struct {
	Id             uuid.UUID        `json:"id"`
	Name           string           `json:"name"`
	RawMaterials   []RawMaterialDTO `json:"raw_materials"`
	Quantity       int              `json:"quantity"`
	Value          MoneyDTO         `json:"value"`
	ProductionCost MoneyDTO         `json:"production_cost"`
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

	rawMaterials := make([]RawMaterialDTO, len(product.RawMaterials()))
	for i, rawMaterial := range product.RawMaterials() {
		rawMaterials[i] = MapRawMaterialEntityToDTO(rawMaterial)
	}

	getMoneyDTO := func(m value_objects.Money) MoneyDTO {
		return MapMoneyObjectToDTO(
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
