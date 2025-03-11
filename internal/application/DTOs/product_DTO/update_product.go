package product_DTO

import (
	"github.com/google/uuid"
	"store-manager/internal/application/DTOs/money_DTO"
)

type UpdateProductDTO struct {
	Id       uuid.UUID          `json:"id"`
	Name     string             `json:"name"`
	Quantity int                `json:"quantity"`
	Value    money_DTO.MoneyDTO `json:"value"`
}
