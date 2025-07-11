package item

import (
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"time"
)

type ItemType string

const (
	Material     ItemType = "MATERIAL"
	Manufactured ItemType = "MANUFACTURED"
)

type Item struct {
	id                string
	name              string
	sku               string
	description       string
	itemType          ItemType
	canBeSold         bool
	active            bool
	priceInCents      int64
	priceCostInCents  int64
	stockQuantity     float64
	minimumStockLevel float64
	unitOfMeasure     string
	deletedAt         *time.Time
}

type HydrateItemParams struct {
	Id                string
	Name              string
	Sku               string
	Description       string
	ItemType          ItemType
	CanBeSold         bool
	Active            bool
	PriceInCents      int64
	PriceCostInCents  int64
	StockQuantity     float64
	MinimumStockLevel float64
	UnitOfMeasure     string
}

func HydrateItem(params HydrateItemParams, deleteAt *time.Time) *Item {
	return &Item{
		id:                params.Id,
		name:              params.Name,
		sku:               params.Sku,
		description:       params.Description,
		itemType:          params.ItemType,
		active:            params.Active,
		canBeSold:         params.CanBeSold,
		priceInCents:      params.PriceInCents,
		priceCostInCents:  params.PriceCostInCents,
		stockQuantity:     params.StockQuantity,
		unitOfMeasure:     params.UnitOfMeasure,
		minimumStockLevel: params.MinimumStockLevel,
		deletedAt:         deleteAt,
	}
}

func (i *Item) ID() string {
	return i.id
}

func (i *Item) Name() string {
	return i.name
}

func (i *Item) SKU() string {
	return i.sku
}

func (i *Item) Description() string {
	return i.description
}

func (i *Item) ItemType() ItemType {
	return i.itemType
}

func (i *Item) CanBeSold() bool {
	return i.canBeSold
}

func (i *Item) SetCanBeSold(canBeSold bool) {
	i.canBeSold = canBeSold
}

func (i *Item) PriceInCents() int64 {
	return i.priceInCents
}

func (i *Item) PriceCostInCents() int64 { return i.priceCostInCents }

func (i *Item) StockQuantity() float64 {
	return i.stockQuantity
}

func (i *Item) UnitOfMeasure() string {
	return i.unitOfMeasure
}

func (i *Item) IsMaterial() bool {
	return i.itemType == Material
}

func (i *Item) IsManufactured() bool {
	return i.itemType == Manufactured
}

func (i *Item) IsAvailableForSale() bool {
	return i.canBeSold && i.stockQuantity > 0
}

func (i *Item) IncreaseStock(amount float64) error {
	if amount < 0 {
		return &domain.ErrInvalidInput{FieldName: "amount", Reason: "amount must be non-negative"}
	}

	i.stockQuantity += amount

	return nil
}

func (i *Item) DecreaseStock(amount float64) error {
	if amount < 0 {
		return &domain.ErrInvalidInput{FieldName: "amount", Reason: "amount must be non-negative"}
	}

	if i.stockQuantity < amount {
		return &domain.ErrInvalidInput{FieldName: "amount", Reason: "insufficient stock"}
	}

	i.stockQuantity -= amount
	return nil
}

func (i *Item) SetPrice(priceInCents int64) error {
	if priceInCents < 0 {
		return &domain.ErrInvalidInput{FieldName: "priceInCents", Reason: "price must be non-negative"}
	}

	i.priceInCents = priceInCents
	return nil
}

func (i *Item) ChangeName(name string) error {
	if name == "" {
		return &domain.ErrInvalidInput{FieldName: "name", Reason: "name is required"}
	}
	i.name = name

	return nil
}

func (i *Item) ChangeDescription(description string) error {
	if description == "" {
		return &domain.ErrInvalidInput{FieldName: "description", Reason: "description is required"}
	}

	i.description = description
	return nil
}

func (i *Item) IsActive() bool {
	return i.active
}

func (i *Item) Activate() {
	i.active = true
}

func (i *Item) Deactivate() {
	i.active = false
}

func (i *Item) MinimumStockLevel() float64 {
	return i.stockQuantity
}

func (i *Item) ChangeMinimumStockLevel(quantity float64) error {
	if quantity <= 0 {
		return &domain.ErrInvalidInput{FieldName: "quantity", Reason: "quantity must be non-negative"}
	}

	i.minimumStockLevel = quantity

	return nil
}

func (i *Item) IsDeleted() bool {
	return i.deletedAt != nil
}

func (i *Item) SetDeleted(deletedAt *time.Time) {
	i.deletedAt = deletedAt
}

func (i *Item) DeletedAt() *time.Time {
	return i.deletedAt
}

func (i *Item) ChangePriceCostInCents(price int64) error {
	if price < 0 {
		return &domain.ErrInvalidInput{FieldName: "Price Cost in Cents", Reason: "price must be non-negative"}
	}

	i.priceCostInCents = price

	return nil
}
