package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Currency         string
	Balance          decimal.Decimal
	AvailableBalance decimal.Decimal
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewAccount(userID uuid.UUID, currency string) *Account {
	now := time.Now()
	return &Account{
		ID:               uuid.New(),
		UserID:           userID,
		Currency:         currency,
		Balance:          decimal.Zero,
		AvailableBalance: decimal.Zero,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (a *Account) HasSufficientBalance(amount decimal.Decimal) bool {
	return a.AvailableBalance.GreaterThanOrEqual(amount)
}

func (a *Account) Debit(amount decimal.Decimal) error {
	if !a.HasSufficientBalance(amount) {
		return ErrInsufficientBalance
	}
	a.Balance = a.Balance.Sub(amount)
	a.AvailableBalance = a.AvailableBalance.Sub(amount)
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) Credit(amount decimal.Decimal) {
	a.Balance = a.Balance.Add(amount)
	a.AvailableBalance = a.AvailableBalance.Add(amount)
	a.UpdatedAt = time.Now()
}

func (a *Account) PlaceHold(amount decimal.Decimal) error {
	if !a.HasSufficientBalance(amount) {
		return ErrInsufficientBalance
	}
	a.AvailableBalance = a.AvailableBalance.Sub(amount)
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) ReleaseHold(amount decimal.Decimal) {
	a.AvailableBalance = a.AvailableBalance.Add(amount)
	a.UpdatedAt = time.Now()
}
