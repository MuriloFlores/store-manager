package value_objects

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrCurrencyMustBeDefined  = errors.New("currency must be defined")
	ErrIncompatibleCurrencies = errors.New("incompatible currencies")
)

type Money struct {
	totalCents int64
	currency   string
}

func NewMoney(dollars, cents int64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, ErrCurrencyMustBeDefined
	}
	totalCents := (dollars * 100) + cents
	return Money{totalCents: totalCents, currency: currency}, nil
}

func (m Money) TotalInDollarsAndCents() (dollars int64, cents int64) {
	return m.Dollars(), m.Cents()
}

func (m Money) ValueInCents() int64 {
	return m.totalCents
}

func (m Money) Dollars() int64 {
	return m.totalCents / 100
}

func (m Money) Cents() int64 {
	return m.totalCents % 100
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrIncompatibleCurrencies
	}
	return Money{totalCents: m.totalCents + other.totalCents, currency: m.currency}, nil
}

func (m Money) Sub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrIncompatibleCurrencies
	}
	return Money{
		totalCents: m.totalCents - other.totalCents,
		currency:   m.currency}, nil
}

func (m Money) Mul(multiplier int64) Money {
	return Money{totalCents: m.totalCents * multiplier, currency: m.currency}
}

func (m Money) Div(divisor int64) Money {
	return Money{totalCents: m.totalCents / divisor, currency: m.currency}
}

func (m Money) String() string {
	return fmt.Sprintf("%s %d.%02d", m.currency, m.Dollars(), m.Cents())
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TotalCents int64  `json:"total_cents"`
		Currency   string `json:"currency"`
	}{
		TotalCents: m.totalCents,
		Currency:   m.currency,
	})
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var aux struct {
		TotalCents int64  `json:"total_cents"`
		Currency   string `json:"currency"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Currency == "" {
		return ErrCurrencyMustBeDefined
	}
	m.totalCents = aux.TotalCents
	m.currency = aux.Currency
	return nil
}
