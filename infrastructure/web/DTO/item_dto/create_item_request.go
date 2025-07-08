package item_dto

type CreateItemRequest struct {
	Name              string  `json:"name" validate:"required,min=3"`
	SKU               string  `json:"SKU" validate:"required,min=3"`
	Description       string  `json:"description" validate:"required"`
	ItemType          string  `json:"item_type" validate:"required,oneof=MANUFACTURED MATERIAL"`
	CanBeSold         bool    `json:"can_be_sold" validate:"required"`
	Active            bool    `json:"active" validate:"required"`
	PriceInCents      int64   `json:"price_in_cents" validate:"required"`
	PriceCostInCents  int64   `json:"price_cost_in_cents" validate:"required"`
	StockQuantity     float64 `json:"stock_quantity" validate:"required"`
	MinimumStockLevel float64 `json:"minimum_stock_level" validate:"required"`
	UnitOfMeasure     string  `json:"unit_of_measure" validate:"required"`
}
