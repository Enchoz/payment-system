package valueobjects

import (
	"errors"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAmount = errors.New("amount must be greater than zero")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

func NewMoney(amount decimal.Decimal, currency string) (*Money, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}
	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m *Money) Add(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, ErrCurrencyMismatch
	}
	return &Money{
		Amount:   m.Amount.Add(other.Amount),
		Currency: m.Currency,
	}, nil
}

func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		return nil, ErrCurrencyMismatch
	}
	if m.Amount.LessThan(other.Amount) {
		return nil, ErrInvalidAmount
	}
	return &Money{
		Amount:   m.Amount.Sub(other.Amount),
		Currency: m.Currency,
	}, nil
}

func (m *Money) IsGreaterThan(other *Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount.GreaterThan(other.Amount)
}

func (m *Money) IsGreaterThanOrEqual(other *Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount.GreaterThanOrEqual(other.Amount)
}
