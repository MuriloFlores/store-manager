package product_DTO

import "github.com/google/uuid"

type FindProductDTO struct {
	Id uuid.UUID `json:"id"`
}
