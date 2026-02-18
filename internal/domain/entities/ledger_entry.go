package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type EntryType string

const (
	EntryTypeDebit  EntryType = "DEBIT"
	EntryTypeCredit EntryType = "CREDIT"
)

type LedgerEntry struct {
	ID          uuid.UUID
	PaymentID   uuid.UUID
	AccountID   uuid.UUID
	EntryType   EntryType
	Amount      decimal.Decimal
	Currency    string
	BalanceAfter decimal.Decimal
	CreatedAt   time.Time
}

func NewLedgerEntry(paymentID, accountID uuid.UUID, entryType EntryType, amount decimal.Decimal, currency string, balanceAfter decimal.Decimal) *LedgerEntry {
	return &LedgerEntry{
		ID:          uuid.New(),
		PaymentID:   paymentID,
		AccountID:   accountID,
		EntryType:   entryType,
		Amount:      amount,
		Currency:    currency,
		BalanceAfter: balanceAfter,
		CreatedAt:   time.Now(),
	}
}
