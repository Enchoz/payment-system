package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentType string

const (
	PaymentTypeInternal PaymentType = "INTERNAL"
	PaymentTypeExternal PaymentType = "EXTERNAL"
)

type PaymentRequestStatus string

const (
	PaymentRequestStatusPending    PaymentRequestStatus = "PENDING"
	PaymentRequestStatusValidating PaymentRequestStatus = "VALIDATING"
	PaymentRequestStatusValidated  PaymentRequestStatus = "VALIDATED"
	PaymentRequestStatusProcessing PaymentRequestStatus = "PROCESSING"
	PaymentRequestStatusCompleted  PaymentRequestStatus = "COMPLETED"
	PaymentRequestStatusFailed     PaymentRequestStatus = "FAILED"
	PaymentRequestStatusCancelled  PaymentRequestStatus = "CANCELLED"
)

type PaymentRequest struct {
	ID                      uuid.UUID
	RequestReference        string
	FromUserID              uuid.UUID
	ExternalAccountID       *uuid.UUID
	ToUserID                *uuid.UUID
	Amount                  decimal.Decimal
	Currency                string
	TargetCurrency          string
	ExternalAccountNumber   string
	ExternalRoutingNumber   string
	ExternalIBAN            string
	ExternalSwiftCode       string
	ExternalBankName        string
	ExternalAccountHolderName string
	ExternalCountryCode     string
	Status                  PaymentRequestStatus
	PaymentType             PaymentType
	ValidationResults       map[string]interface{}
	ValidationErrors        map[string]interface{}
	RequiresManualReview    bool
	ExchangeRate            *decimal.Decimal
	ConvertedAmount         *decimal.Decimal
	FailureReason           string
	ProcessedAt             *time.Time
	Metadata                map[string]interface{}
	IdempotencyKey          string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func NewInternalPaymentRequest(fromUserID, toUserID uuid.UUID, amount decimal.Decimal, currency string) *PaymentRequest {
	now := time.Now()
	return &PaymentRequest{
		ID:           uuid.New(),
		FromUserID:   fromUserID,
		ToUserID:     &toUserID,
		Amount:       amount,
		Currency:     currency,
		PaymentType:  PaymentTypeInternal,
		Status:       PaymentRequestStatusPending,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func NewExternalPaymentRequest(fromUserID uuid.UUID, amount decimal.Decimal, currency, targetCurrency string) *PaymentRequest {
	now := time.Now()
	return &PaymentRequest{
		ID:            uuid.New(),
		FromUserID:    fromUserID,
		Amount:        amount,
		Currency:      currency,
		TargetCurrency: targetCurrency,
		PaymentType:   PaymentTypeExternal,
		Status:        PaymentRequestStatusPending,
		Metadata:      make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (pr *PaymentRequest) MarkAsValidating() {
	pr.Status = PaymentRequestStatusValidating
	pr.UpdatedAt = time.Now()
}

func (pr *PaymentRequest) MarkAsValidated() {
	pr.Status = PaymentRequestStatusValidated
	pr.UpdatedAt = time.Now()
}

func (pr *PaymentRequest) MarkAsProcessing() {
	pr.Status = PaymentRequestStatusProcessing
	now := time.Now()
	pr.ProcessedAt = &now
	pr.UpdatedAt = now
}

func (pr *PaymentRequest) MarkAsCompleted() {
	pr.Status = PaymentRequestStatusCompleted
	pr.UpdatedAt = time.Now()
}

func (pr *PaymentRequest) MarkAsFailed(reason string) {
	pr.Status = PaymentRequestStatusFailed
	pr.FailureReason = reason
	pr.UpdatedAt = time.Now()
}

func (pr *PaymentRequest) SetValidationResults(results map[string]interface{}) {
	pr.ValidationResults = results
	pr.UpdatedAt = time.Now()
}

func (pr *PaymentRequest) SetValidationErrors(errors map[string]interface{}) {
	pr.ValidationErrors = errors
	pr.UpdatedAt = time.Now()
}

func (pr *PaymentRequest) SetExchangeRate(rate decimal.Decimal, convertedAmount decimal.Decimal) {
	pr.ExchangeRate = &rate
	pr.ConvertedAmount = &convertedAmount
	pr.UpdatedAt = time.Now()
}
