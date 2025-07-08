package item_dto

import "github.com/muriloFlores/StoreManager/internal/core/domain/item"

type ClientItemResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SKU         string `json:"sku"`
	PriceSale   int64  `json:"price_sale"`
	Active      bool   `json:"active"`
}

func ToClientItemResponse(item *item.Item) ClientItemResponse {
	return ClientItemResponse{
		ID:          item.ID(),
		Name:        item.Name(),
		Description: item.Description(),
		SKU:         item.SKU(),
		PriceSale:   item.PriceInCents(),
		Active:      item.IsActive(),
	}
}
