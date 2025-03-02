package entity

import (
	"sort"
	"store-manager/internal/domain/value_objects"
)

type product struct {
	id             EntityID
	name           string
	rawMaterials   []RawMaterialInterface
	quantity       int
	value          value_objects.Money
	productionCost value_objects.Money
}

type ProductInterface interface {
	ID() EntityID
	Name() string
	RawMaterials() []RawMaterialInterface
	Quantity() int
	Value() value_objects.Money
	ProductionCost() value_objects.Money

	SetValue(value value_objects.Money)
	SetQuantity(quantity int)
	SetName(name string)

	AddRawMaterials(rawMaterials []RawMaterialInterface)
	RemoveRawMaterials(rawMaterials []RawMaterialInterface)

	CalculateCost()
	OrganizeRawMaterialsByCost()
	OrganizeRawMaterialsByQuantity()
	OrganizeRawMaterialsByName()
}

func NewProduct(id *EntityID, name string, rawMaterials []RawMaterialInterface, quantity int, value value_objects.Money) ProductInterface {
	var productId EntityID
	if id == nil {
		productId = NewEntityID()
	} else {
		productId = *id
	}

	newProduct := product{
		id:           productId,
		name:         name,
		rawMaterials: rawMaterials,
		quantity:     quantity,
		value:        value,
	}

	newProduct.CalculateCost()

	return &newProduct
}

func (p *product) ID() EntityID {
	return p.id
}

func (p *product) Name() string {
	return p.name
}

func (p *product) RawMaterials() []RawMaterialInterface {
	return p.rawMaterials
}

func (p *product) Quantity() int {
	return p.quantity
}

func (p *product) Value() value_objects.Money {
	return p.value
}

func (p *product) ProductionCost() value_objects.Money {
	return p.productionCost
}

func (p *product) SetValue(value value_objects.Money) {
	p.value = value
}

func (p *product) SetQuantity(quantity int) {
	p.quantity = quantity
}

func (p *product) SetName(name string) {
	p.name = name
}

func (p *product) AddRawMaterials(rawMaterials []RawMaterialInterface) {
	p.rawMaterials = append(p.rawMaterials, rawMaterials...)
}

func (p *product) RemoveRawMaterials(RemoveItems []RawMaterialInterface) {
	var filteredMaterials []RawMaterialInterface

	for _, material := range p.rawMaterials {
		shouldRemove := false
		for _, removeItem := range RemoveItems {
			if material.Name() == removeItem.Name() {
				shouldRemove = true
				break
			}
		}

		if !shouldRemove {
			filteredMaterials = append(filteredMaterials, material)
		}
	}

	p.rawMaterials = filteredMaterials
}

func (p *product) OrganizeRawMaterialsByCost() {
	sort.Slice(p.rawMaterials, func(i, j int) bool {
		return p.rawMaterials[i].Cost().ValueInCents() < p.rawMaterials[j].Cost().ValueInCents()
	})
}

func (p *product) OrganizeRawMaterialsByQuantity() {
	sort.Slice(p.rawMaterials, func(i, j int) bool {
		return p.rawMaterials[i].Quantity() < p.rawMaterials[j].Quantity()
	})
}

func (p *product) OrganizeRawMaterialsByName() {
	sort.Slice(p.rawMaterials, func(i, j int) bool {
		return p.rawMaterials[i].Name() < p.rawMaterials[j].Name()
	})
}

func (p *product) CalculateCost() {
	if len(p.rawMaterials) == 0 {
		none, _ := value_objects.NewMoney(0, 0, p.value.Currency())

		p.productionCost = none
		return
	}

	baseCurrency := p.rawMaterials[0].Cost().Currency()

	sum, _ := value_objects.NewMoney(0, 0, p.value.Currency())

	for _, material := range p.rawMaterials {
		cost := material.Cost()
		sum, _ = sum.Add(cost)
	}

	totalCents := sum.ValueInCents()
	laborCostInCents := (totalCents * 10) / 100

	laborDollars := laborCostInCents / 100
	laborCents := laborCostInCents % 100

	laborCost, _ := value_objects.NewMoney(laborDollars, laborCents, baseCurrency)
	finalCost, _ := sum.Add(laborCost)

	p.productionCost = finalCost
}
