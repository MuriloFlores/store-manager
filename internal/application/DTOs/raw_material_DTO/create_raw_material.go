package raw_material_DTO

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"store-manager/internal/application/DTOs/money_DTO"
	"store-manager/internal/domain/entity"
)

type CreateRawMaterialDTO struct {
	Name      string             `json:"name"`
	Unit      entity.Unit        `json:"unit"`
	Quantity  int                `json:"quantity"`
	Cost      money_DTO.MoneyDTO `json:"cost"`
	RiskLimit int                `json:"risk_limit"`
}

func (cr *CreateRawMaterialDTO) Validate() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.Name, validation.Required),
		validation.Field(&cr.Unit, validation.Required),
		validation.Field(&cr.Quantity, validation.Required, validation.Min(0)),
		validation.Field(&cr.Cost, validation.Required),
	)
}
