package money_DTO

import (
	"store-manager/internal/domain/value_objects"
)

type MoneyDTO struct {
	TotalCents   int64  `json:"total_cents"`
	CurrencyCode string `json:"currency_code"`
}

func (m *MoneyDTO) MapMoneyDTOToObject() value_objects.Money {
	money, _ := value_objects.NewMoney(
		m.TotalCents/100, //dollar
		m.TotalCents%100, //cents
		m.CurrencyCode,
	)

	return money
}

func MapMoneyObjectToDTO(money value_objects.Money) MoneyDTO {
	return MoneyDTO{
		money.ValueInCents(),
		money.Currency(),
	}
}
