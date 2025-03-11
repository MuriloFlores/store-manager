package product_DTO

import "store-manager/internal/application/DTOs/money_DTO"

type CreateProductDTO struct {
	Name     string             `json:"name"`
	Quantity int                `json:"quantity"`
	Value    money_DTO.MoneyDTO `json:"value"`
}
