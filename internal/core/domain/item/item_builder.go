package item

import (
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"time"
)

type ItemBuilder struct {
	item Item
	err  error
}

func NewItemBuilder() *ItemBuilder {
	return &ItemBuilder{}
}

func (b *ItemBuilder) WithID(id string) *ItemBuilder {
	if id == "" {
		b.err = &domain.ErrInvalidInput{FieldName: "id", Reason: "id is required"}
	}

	b.item.id = id

	return b
}

func (b *ItemBuilder) WithName(name string) *ItemBuilder {
	if name == "" {
		b.err = &domain.ErrInvalidInput{FieldName: "name", Reason: "name is required"}
	}

	b.item.name = name

	return b
}

func (b *ItemBuilder) WithSKU(sku string) *ItemBuilder {
	b.item.sku = sku
	return b
}

func (b *ItemBuilder) WithDescription(description string) *ItemBuilder {
	b.item.description = description

	return b
}

func (b *ItemBuilder) WithType(itemType ItemType) *ItemBuilder {
	b.item.itemType = itemType

	return b
}

func (b *ItemBuilder) WithPriceInCents(price int64) *ItemBuilder {
	if price < 0 {
		b.err = &domain.ErrInvalidInput{FieldName: "priceInCents", Reason: "price must be non-negative"}
	}

	b.item.priceInCents = price

	return b
}

func (b *ItemBuilder) WithUnitOfMeasure(measure string) *ItemBuilder {
	if measure == "" {
		b.err = &domain.ErrInvalidInput{FieldName: "measure", Reason: "measure is required"}
	}

	b.item.unitOfMeasure = measure

	return b
}

func (b *ItemBuilder) WithQuantity(quantity float64) *ItemBuilder {
	if quantity < 0 {
		b.err = &domain.ErrInvalidInput{FieldName: "quantity", Reason: "quantity must be non-negative"}
	}

	b.item.stockQuantity = quantity

	return b
}

func (b *ItemBuilder) WithCanBeSold(beSold bool) *ItemBuilder {
	b.item.canBeSold = beSold

	return b
}

func (b *ItemBuilder) WithMinimumStock(minimum float64) *ItemBuilder {
	if minimum < 0 {
		b.err = &domain.ErrInvalidInput{FieldName: "minimumStock", Reason: "minimum stock must be non-negative"}
	}

	b.item.minimumStockLevel = minimum

	return b
}

func (b *ItemBuilder) WithDeletion(date *time.Time) *ItemBuilder {
	b.item.deletedAt = date

	return b
}

func (b *ItemBuilder) WithPriceCostInCents(price int64) *ItemBuilder {
	if price < 0 {
		b.err = &domain.ErrInvalidInput{FieldName: "priceCostInCents", Reason: "price must be non-negative"}
	}

	b.item.priceCostInCents = price

	return b
}

func (b *ItemBuilder) WithActive(active bool) *ItemBuilder {
	b.item.active = active

	return b
}

func (b *ItemBuilder) Build() (*Item, error) {
	if b.err != nil {
		return nil, b.err
	}

	if b.item.itemType == "" {
		b.item.itemType = Material
	}

	b.item.active = true

	return &b.item, nil
}
