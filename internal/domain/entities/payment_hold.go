package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentHold struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	RequestID  *uuid.UUID
	PaymentID  *uuid.UUID
	Amount     decimal.Decimal
	Currency   string
	ExpiresAt  *time.Time
	ReleasedAt *time.Time
	CreatedAt  time.Time
}

func NewPaymentHold(accountID uuid.UUID, amount decimal.Decimal, currency string, requestID *uuid.UUID) *PaymentHold {
	return &PaymentHold{
		ID:        uuid.New(),
		AccountID: accountID,
		RequestID: requestID,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: time.Now(),
	}
}

func (ph *PaymentHold) Release() {
	now := time.Now()
	ph.ReleasedAt = &now
}

func (ph *PaymentHold) IsReleased() bool {
	return ph.ReleasedAt != nil
}

func (ph *PaymentHold) IsExpired() bool {
	if ph.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ph.ExpiresAt)
}
