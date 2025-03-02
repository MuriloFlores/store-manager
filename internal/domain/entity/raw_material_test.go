package entity

import (
	"github.com/stretchr/testify/assert"
	"store-manager/internal/domain/value_objects"
	"testing"
)

func TestRawMaterial_NewRaMaterial(t *testing.T) {
	cost, err := value_objects.NewMoney(10, 50, "USD")
	assert.NoError(t, err)

	rm := NewRawMaterial(nil, "Steel", "kg", 100, cost)

	assert.NotNil(t, rm)
	assert.Equal(t, "Steel", rm.Name())

	assert.Equal(t, "kg", string(rm.Unit()))
	assert.Equal(t, 100, rm.Quantity())
	assert.Equal(t, cost.String(), rm.Cost().String())
}

func TestRawMaterial_SetName(t *testing.T) {
	cost, err := value_objects.NewMoney(10, 50, "USD")
	assert.NoError(t, err)

	rm := NewRawMaterial(nil, "Steel", "kg", 100, cost)

	rm.SetName("Aluminum")
	assert.Equal(t, "Aluminum", rm.Name(), "Expected name to be updated to 'Aluminum'")
}

func TestRawMaterial_SetUnit(t *testing.T) {
	cost, err := value_objects.NewMoney(10, 50, "USD")
	assert.NoError(t, err)

	rm := NewRawMaterial(nil, "Steel", "kg", 100, cost)

	err = rm.SetUnit("lts")
	assert.NoError(t, err)
	assert.Equal(t, "lts", string(rm.Unit()), "Expected unit to be updated to 'lts'")

	err = rm.SetUnit("invalid")
	assert.Error(t, err, "Expected error when setting an invalid unit")

	assert.Equal(t, "lts", string(rm.Unit()), "Unit should remain unchanged when setting an invalid value")
}

func TestRawMaterial_SetQuantity(t *testing.T) {
	cost, err := value_objects.NewMoney(10, 50, "USD")
	assert.NoError(t, err)

	rm := NewRawMaterial(nil, "Steel", "kg", 100, cost)
	rm.SetQuantity(200)

	assert.Equal(t, 200, rm.Quantity(), "Expected quantity to be updated to 200")
}

func TestRawMaterial_SetCost(t *testing.T) {
	initialCost, err := value_objects.NewMoney(10, 50, "USD")
	assert.NoError(t, err)

	newCost, err := value_objects.NewMoney(20, 00, "USD")
	assert.NoError(t, err)

	rm := NewRawMaterial(nil, "Steel", "kg", 100, initialCost)
	rm.SetCost(newCost)

	assert.Equal(t, newCost.String(), rm.Cost().String(), "Expected cost to be updated")
}
