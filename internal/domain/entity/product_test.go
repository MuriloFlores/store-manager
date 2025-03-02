package entity

import (
	"github.com/stretchr/testify/assert"
	"store-manager/internal/domain/value_objects"
	"testing"
)

func TestProduct_OrganizeRawMaterialsByCost(t *testing.T) {
	m5, err := value_objects.NewMoney(5, 0, "USD")
	assert.NoError(t, err)

	m10, err := value_objects.NewMoney(10, 0, "USD")
	assert.NoError(t, err)

	m15, err := value_objects.NewMoney(15, 0, "USD")
	assert.NoError(t, err)

	rm1 := NewRawMaterial(nil, "Material1", "kg", 5, m10)
	rm2 := NewRawMaterial(nil, "Material2", "kg", 5, m5)
	rm3 := NewRawMaterial(nil, "Material3", "kg", 5, m15)

	rawMaterials := []RawMaterialInterface{rm1, rm2, rm3}
	prod := NewProduct(nil, "product1", rawMaterials, 100, m15)

	prod.OrganizeRawMaterialsByCost()
	sorted := prod.RawMaterials()

	getTotalInCents := func(money value_objects.Money) int64 {
		return money.ValueInCents()
	}

	cost0 := getTotalInCents(sorted[0].Cost())
	cost1 := getTotalInCents(sorted[1].Cost())
	cost2 := getTotalInCents(sorted[2].Cost())

	assert.True(t, cost0 <= cost1 && cost1 <= cost2, "Raw materials are not sorted by cost")
}

func TestProduct_OrganizeRawMaterialsByQuantity(t *testing.T) {
	money, _ := value_objects.NewMoney(1, 0, "USD")

	rm1 := NewRawMaterial(nil, "Material1", "kg", 20, money)
	rm2 := NewRawMaterial(nil, "Material2", "kg", 5, money)
	rm3 := NewRawMaterial(nil, "Material3", "kg", 10, money)

	rawMaterials := []RawMaterialInterface{rm1, rm2, rm3}
	prod := NewProduct(nil, "product1", rawMaterials, 50, money)

	prod.OrganizeRawMaterialsByQuantity()
	sorted := prod.RawMaterials()

	assert.Equal(t, 5, sorted[0].Quantity(), "First raw material should have quantity 5")
	assert.Equal(t, 10, sorted[1].Quantity(), "First raw material should have quantity 10")
	assert.Equal(t, 20, sorted[2].Quantity(), "First raw material should have quantity 20")
}

func TestProduct_OrganizeRawMaterialsByName(t *testing.T) {
	money, _ := value_objects.NewMoney(1, 0, "USD")

	rm1 := NewRawMaterial(nil, "zeta", "kg", 20, money)
	rm2 := NewRawMaterial(nil, "alpha", "kg", 5, money)
	rm3 := NewRawMaterial(nil, "beta", "kg", 10, money)

	rawMaterials := []RawMaterialInterface{rm1, rm2, rm3}

	prod := NewProduct(nil, "product1", rawMaterials, 50, money)

	prod.OrganizeRawMaterialsByName()
	sorted := prod.RawMaterials()

	assert.Equal(t, "alpha", sorted[0].Name())
	assert.Equal(t, "beta", sorted[1].Name())
	assert.Equal(t, "zeta", sorted[2].Name())
}

func TestProduct_AddRawMaterials(t *testing.T) {
	money, _ := value_objects.NewMoney(1, 0, "USD")
	prod := NewProduct(nil, "Product", []RawMaterialInterface{}, 50, money)

	rm := NewRawMaterial(nil, "Product", "kg", 10, money)

	prod.AddRawMaterials([]RawMaterialInterface{rm})
	assert.Equal(t, 1, len(prod.RawMaterials()), "Raw material was not added")
}

func TestProduct_RemoveRawMaterials(t *testing.T) {
	money, _ := value_objects.NewMoney(1, 0, "USD")

	rm1 := NewRawMaterial(nil, "Material1", "kg", 10, money)
	rm2 := NewRawMaterial(nil, "Material2", "kg", 20, money)
	rm3 := NewRawMaterial(nil, "Material3", "kg", 30, money)

	rawMaterials := []RawMaterialInterface{rm1, rm2, rm3}
	prod := NewProduct(nil, "product1", rawMaterials, 50, money)

	prod.RemoveRawMaterials([]RawMaterialInterface{rm2})
	remaining := prod.RawMaterials()

	for _, rm := range remaining {
		assert.NotEqual(t, "Material2", rm.Name(), "Raw material Material2 was not removed")
	}

	assert.Equal(t, 2, len(remaining), "Number of raw materials is not correct after removal")
}
