package raw_material_DTO

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type FindRawMaterialDTO struct {
	Id uuid.UUID `json:"id"`
}

func (r *FindRawMaterialDTO) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Id, validation.Required),
	)
}
