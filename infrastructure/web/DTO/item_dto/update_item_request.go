package item_dto

type UpdateItemRequest struct {
	Name              *string `json:"name" validate:"omitempty,min=3"`
	Description       *string `json:"description" validate:"omitempty,min=1"`
	IsActive          *bool   `json:"is_active" validate:"omitempty"`
	CanBeSold         *bool   `json:"can_be_sold" validate:"omitempty"`
	PriceSaleInCents  *int64  `json:"price_sale_in_cents" validate:"omitempty,gte=0"`
	MinimumStockLevel *int    `json:"minimum_stock_level" validate:"omitempty,gte=0"`
}
