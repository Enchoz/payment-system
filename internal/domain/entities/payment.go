package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusReversed   PaymentStatus = "REVERSED"
)

type Payment struct {
	ID                    uuid.UUID
	PaymentReference      string
	RequestID             *uuid.UUID
	FromUserID            uuid.UUID
	ToUserID              *uuid.UUID
	ExternalAccountID     *uuid.UUID
	Amount                decimal.Decimal
	Currency              string
	ExchangeRate          *decimal.Decimal
	ConvertedAmount       *decimal.Decimal
	ConvertedCurrency     string
	Status                PaymentStatus
	PaymentType           PaymentType
	ExternalAccountNumber string
	ExternalBankName      string
	FailureReason         string
	ProcessedAt           *time.Time
	CompletedAt           *time.Time
	Metadata              map[string]interface{}
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func NewPayment(requestID *uuid.UUID, fromUserID uuid.UUID, amount decimal.Decimal, currency string, paymentType PaymentType) *Payment {
	now := time.Now()
	return &Payment{
		ID:          uuid.New(),
		RequestID:   requestID,
		FromUserID:  fromUserID,
		Amount:      amount,
		Currency:    currency,
		PaymentType: paymentType,
		Status:      PaymentStatusProcessing,
		ProcessedAt: &now,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Payment) MarkAsCompleted() {
	now := time.Now()
	p.Status = PaymentStatusCompleted
	p.CompletedAt = &now
	p.UpdatedAt = now
}

func (p *Payment) MarkAsFailed(reason string) {
	p.Status = PaymentStatusFailed
	p.FailureReason = reason
	p.UpdatedAt = time.Now()
}

func (p *Payment) SetExchangeDetails(rate decimal.Decimal, convertedAmount decimal.Decimal, convertedCurrency string) {
	p.ExchangeRate = &rate
	p.ConvertedAmount = &convertedAmount
	p.ConvertedCurrency = convertedCurrency
	p.UpdatedAt = time.Now()
}
