package item_dto

import "github.com/muriloFlores/StoreManager/internal/core/domain/item"

type InternalItemResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	SKU               string  `json:"sku"`
	ItemType          string  `json:"item_type"`
	Active            bool    `json:"active"`
	CanBeSold         bool    `json:"can_be_sold"`
	PriceSale         int64   `json:"price_sale"`
	PriceCost         int64   `json:"price_cost"`
	StockQuantity     float64 `json:"stock_quantity"`
	UnitOfMeasure     string  `json:"unit_of_measure"`
	MinimumStockLevel float64 `json:"minimum_stock_level"`
}

func ToInternalItemResponse(item *item.Item) InternalItemResponse {
	return InternalItemResponse{
		ID:                item.ID(),
		Name:              item.Name(),
		Description:       item.Description(),
		SKU:               item.SKU(),
		ItemType:          string(item.ItemType()),
		Active:            item.IsActive(),
		CanBeSold:         item.CanBeSold(),
		PriceSale:         item.PriceInCents(),
		PriceCost:         item.PriceCostInCents(),
		StockQuantity:     item.StockQuantity(),
		UnitOfMeasure:     item.UnitOfMeasure(),
		MinimumStockLevel: item.MinimumStockLevel(),
	}
}
