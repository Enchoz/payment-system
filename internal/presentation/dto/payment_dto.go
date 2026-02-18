package dto

import (
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InternalPaymentRequest struct {
	FromUserID uuid.UUID       `json:"fromUserId" binding:"required"`
	ToUserID   uuid.UUID       `json:"toUserId" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required"`
	Currency   string          `json:"currency" binding:"required"`
	Reference  string          `json:"reference" binding:"required"`
}

type ExternalPaymentRequest struct {
	FromUserID                uuid.UUID  `json:"fromUserId" binding:"required"`
	ExternalAccountID         *uuid.UUID `json:"externalAccountId"`
	Amount                    string     `json:"amount" binding:"required"`
	Currency                  string     `json:"currency" binding:"required"`
	TargetCurrency            string     `json:"targetCurrency"`
	ExternalAccountNumber     string     `json:"externalAccountNumber"`
	ExternalRoutingNumber     string     `json:"externalRoutingNumber"`
	ExternalIBAN              string     `json:"externalIban"`
	ExternalSwiftCode         string     `json:"externalSwiftCode"`
	ExternalBankName          string     `json:"externalBankName"`
	ExternalAccountHolderName string     `json:"externalAccountHolderName"`
	ExternalCountryCode       string     `json:"externalCountryCode"`
	Reference                 string     `json:"reference" binding:"required"`
	IdempotencyKey            string     `json:"idempotencyKey"`
	SaveAccount               bool       `json:"saveAccount"`
}

type PaymentResponse struct {
	ID                uuid.UUID  `json:"id"`
	PaymentReference  string     `json:"paymentReference"`
	RequestID         *uuid.UUID `json:"requestId"`
	FromUserID        uuid.UUID  `json:"fromUserId"`
	ToUserID          *uuid.UUID `json:"toUserId"`
	ExternalAccountID *uuid.UUID `json:"externalAccountId"`
	Amount            string     `json:"amount"`
	Currency          string     `json:"currency"`
	ExchangeRate       *string   `json:"exchangeRate"`
	ConvertedAmount   *string   `json:"convertedAmount"`
	ConvertedCurrency string     `json:"convertedCurrency"`
	Status            string     `json:"status"`
	PaymentType       string     `json:"paymentType"`
	CreatedAt         string     `json:"createdAt"`
	CompletedAt       *string    `json:"completedAt"`
}

type PaymentRequestResponse struct {
	ID                  uuid.UUID              `json:"id"`
	RequestReference    string                 `json:"requestReference"`
	FromUserID          uuid.UUID              `json:"fromUserId"`
	ExternalAccountID   *uuid.UUID             `json:"externalAccountId"`
	ToUserID            *uuid.UUID             `json:"toUserId"`
	Amount              string                 `json:"amount"`
	Currency            string                 `json:"currency"`
	TargetCurrency      string                 `json:"targetCurrency"`
	Status              string                 `json:"status"`
	PaymentType         string                 `json:"paymentType"`
	ValidationResults   map[string]interface{} `json:"validationResults"`
	ValidationErrors    map[string]interface{} `json:"validationErrors"`
	ExchangeRate        *string                `json:"exchangeRate"`
	ConvertedAmount     *string                `json:"convertedAmount"`
	RequiresManualReview bool                  `json:"requiresManualReview"`
	CreatedAt           string                 `json:"createdAt"`
	UpdatedAt           string                 `json:"updatedAt"`
}

func PaymentToResponse(payment *entities.Payment) *PaymentResponse {
	response := &PaymentResponse{
		ID:                payment.ID,
		PaymentReference:  payment.PaymentReference,
		RequestID:         payment.RequestID,
		FromUserID:        payment.FromUserID,
		ToUserID:          payment.ToUserID,
		ExternalAccountID: payment.ExternalAccountID,
		Amount:            payment.Amount.String(),
		Currency:          payment.Currency,
		Status:            string(payment.Status),
		PaymentType:       string(payment.PaymentType),
		CreatedAt:         payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ConvertedCurrency: payment.ConvertedCurrency,
	}

	if payment.ExchangeRate != nil {
		rate := payment.ExchangeRate.String()
		response.ExchangeRate = &rate
	}

	if payment.ConvertedAmount != nil {
		amount := payment.ConvertedAmount.String()
		response.ConvertedAmount = &amount
	}

	if payment.CompletedAt != nil {
		completedAt := payment.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAt
	}

	return response
}

func PaymentRequestToResponse(request *entities.PaymentRequest) *PaymentRequestResponse {
	response := &PaymentRequestResponse{
		ID:                  request.ID,
		RequestReference:    request.RequestReference,
		FromUserID:          request.FromUserID,
		ExternalAccountID:   request.ExternalAccountID,
		ToUserID:            request.ToUserID,
		Amount:              request.Amount.String(),
		Currency:            request.Currency,
		TargetCurrency:      request.TargetCurrency,
		Status:              string(request.Status),
		PaymentType:         string(request.PaymentType),
		ValidationResults:   request.ValidationResults,
		ValidationErrors:    request.ValidationErrors,
		RequiresManualReview: request.RequiresManualReview,
		CreatedAt:           request.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:           request.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if request.ExchangeRate != nil {
		rate := request.ExchangeRate.String()
		response.ExchangeRate = &rate
	}

	if request.ConvertedAmount != nil {
		amount := request.ConvertedAmount.String()
		response.ConvertedAmount = &amount
	}

	return response
}
