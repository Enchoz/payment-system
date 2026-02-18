package entities

import (
	"time"

	"github.com/google/uuid"
)

type ValidationStatus string

const (
	ValidationStatusPending   ValidationStatus = "PENDING"
	ValidationStatusValidating ValidationStatus = "VALIDATING"
	ValidationStatusValid      ValidationStatus = "VALID"
	ValidationStatusInvalid    ValidationStatus = "INVALID"
	ValidationStatusExpired    ValidationStatus = "EXPIRED"
)

type ExternalAccount struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	AccountNumber     string
	RoutingNumber     string
	IBAN              string
	SwiftCode         string
	BankName          string
	AccountHolderName string
	Currency          string
	AccountType       string
	CountryCode       string
	ValidationStatus  ValidationStatus
	ValidatedAt       *time.Time
	ValidationExpiresAt *time.Time
	LastValidatedAt   *time.Time
	Nickname          string
	IsDefault         bool
	Metadata          map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewExternalAccount(userID uuid.UUID, accountNumber, currency string) *ExternalAccount {
	now := time.Now()
	return &ExternalAccount{
		ID:               uuid.New(),
		UserID:           userID,
		AccountNumber:    accountNumber,
		Currency:         currency,
		ValidationStatus: ValidationStatusPending,
		Metadata:         make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (ea *ExternalAccount) IsValid() bool {
	return ea.ValidationStatus == ValidationStatusValid
}

func (ea *ExternalAccount) NeedsRevalidation() bool {
	if ea.ValidationStatus == ValidationStatusInvalid || ea.ValidationStatus == ValidationStatusExpired {
		return true
	}
	if ea.ValidationExpiresAt != nil && time.Now().After(*ea.ValidationExpiresAt) {
		return true
	}
	return false
}

func (ea *ExternalAccount) MarkAsValid(expiresAt *time.Time) {
	now := time.Now()
	ea.ValidationStatus = ValidationStatusValid
	ea.ValidatedAt = &now
	ea.LastValidatedAt = &now
	ea.ValidationExpiresAt = expiresAt
	ea.UpdatedAt = now
}

func (ea *ExternalAccount) MarkAsInvalid() {
	now := time.Now()
	ea.ValidationStatus = ValidationStatusInvalid
	ea.UpdatedAt = now
}
