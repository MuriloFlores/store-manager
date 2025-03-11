package raw_material_DTO

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"store-manager/internal/application/DTOs/money_DTO"
	"store-manager/internal/domain/entity"
)

type RawMaterialDTO struct {
	Id        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Unit      entity.Unit        `json:"unit"`
	Quantity  int                `json:"quantity"`
	Cost      money_DTO.MoneyDTO `json:"cost"`
	RiskLimit int                `json:"risk_limit"`
}

func (r *RawMaterialDTO) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Id, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Unit, validation.Required),
		validation.Field(&r.Quantity, validation.Required, validation.Min(0)),
		validation.Field(&r.Cost, validation.Required),
	)
}

func (r *RawMaterialDTO) MapRawMaterialDTOToEntity() entity.RawMaterialInterface {
	id, _ := entity.ParseEntityID(r.Id.String())

	return entity.NewRawMaterial(
		&id,
		r.Name,
		r.Unit,
		r.Quantity,
		r.Cost.MapMoneyDTOToObject(),
		&r.RiskLimit,
	)
}

func MapRawMaterialEntityToDTO(e entity.RawMaterialInterface) RawMaterialDTO {
	id, _ := uuid.Parse(e.ID().String())

	totalCostInCentsDTO := money_DTO.MapMoneyObjectToDTO(
		e.Cost(),
	)

	return RawMaterialDTO{
		id,
		e.Name(),
		e.Unit(),
		e.Quantity(),
		totalCostInCentsDTO,
		e.RiskLimit(),
	}
}
