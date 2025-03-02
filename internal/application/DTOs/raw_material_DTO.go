package DTOs

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"store-manager/internal/domain/entity"
)

type RawMaterialDTO struct {
	Id       uuid.UUID   `json:"id"`
	Name     string      `json:"name"`
	Unit     entity.Unit `json:"unit"`
	Quantity int         `json:"quantity"`
	Cost     MoneyDTO    `json:"cost"`
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
	)
}

func MapRawMaterialEntityToDTO(e entity.RawMaterialInterface) RawMaterialDTO {
	id, _ := uuid.Parse(e.ID().String())

	totalCostInCentsDTO := MapMoneyObjectToDTO(
		e.Cost(),
	)

	return RawMaterialDTO{
		id,
		e.Name(),
		e.Unit(),
		e.Quantity(),
		totalCostInCentsDTO,
	}
}

type CreateRawMaterialDTO struct {
	Name     string      `json:"name"`
	Unit     entity.Unit `json:"unit"`
	Quantity int         `json:"quantity"`
	Cost     MoneyDTO    `json:"cost"`
}

func (cr *CreateRawMaterialDTO) Validate() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.Name, validation.Required),
		validation.Field(&cr.Unit, validation.Required),
		validation.Field(&cr.Quantity, validation.Required, validation.Min(0)),
		validation.Field(&cr.Cost, validation.Required),
	)
}

type ReturnRawMaterialDTO struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
}

type FindRawMaterialDTO struct {
	Id uuid.UUID `json:"id"`
}

func (r *FindRawMaterialDTO) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Id, validation.Required),
	)
}
