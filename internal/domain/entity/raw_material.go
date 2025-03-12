package entity

import "store-manager/internal/domain/value_objects"

type rawMaterial struct {
	id        EntityID
	name      string
	unit      Unit
	quantity  int
	riskLimit int
	cost      value_objects.Money
}

type RawMaterialInterface interface {
	ID() EntityID
	Name() string
	Unit() Unit
	Quantity() int
	Cost() value_objects.Money
	RiskLimit() int

	SetName(string)
	SetUnit(Unit) error
	SetQuantity(int)
	SetCost(value_objects.Money)
	SetRiskLimit(int)
}

func NewRawMaterial(id *EntityID, name string, unit Unit, quantity int, cost value_objects.Money, riskLimit *int) RawMaterialInterface {
	var productId EntityID
	if id == nil {
		productId = NewEntityID()
	} else {
		productId = *id
	}

	var computedRiskLimit int
	if riskLimit == nil {
		computedRiskLimit = quantity / 2
	} else {
		computedRiskLimit = *riskLimit
	}

	return &rawMaterial{
		id:        productId,
		name:      name,
		unit:      unit,
		quantity:  quantity,
		cost:      cost,
		riskLimit: computedRiskLimit,
	}
}

func (m *rawMaterial) ID() EntityID {
	return m.id
}

func (m *rawMaterial) Name() string {
	return m.name
}

func (m *rawMaterial) Unit() Unit {
	return m.unit
}

func (m *rawMaterial) Quantity() int {
	return m.quantity
}

func (m *rawMaterial) Cost() value_objects.Money {
	return m.cost
}

func (m *rawMaterial) RiskLimit() int {
	return m.riskLimit
}

func (m *rawMaterial) SetName(name string) {
	m.name = name
}

func (m *rawMaterial) SetUnit(unit Unit) error {
	if err := ValidateUnit(unit); err != nil {
		return err
	}

	m.unit = unit

	return nil
}

func (m *rawMaterial) SetQuantity(quantity int) {
	m.quantity = quantity
}

func (m *rawMaterial) SetCost(cost value_objects.Money) {
	m.cost = cost
}

func (m *rawMaterial) SetRiskLimit(riskLimit int) {
	m.riskLimit = riskLimit
}
