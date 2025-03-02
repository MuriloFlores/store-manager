package models

import "store-manager/internal/domain/value_objects"

type MoneyModel struct {
	TotalInCents int64  `gorm:"not null"`
	Currency     string `gorm:"type:varchar(10);not null"`
}

func (m *MoneyModel) MapMoneyModelToEntity() value_objects.Money {
	money, _ := value_objects.NewMoney(
		m.TotalInCents/100,
		m.TotalInCents%100,
		m.Currency,
	)
	return money
}

func MapMoneyObjectToMoneyModel(money value_objects.Money) MoneyModel {
	return MoneyModel{
		TotalInCents: money.ValueInCents(),
		Currency:     money.Currency(),
	}
}
