package services

import (
	"context"
	"payment-system/internal/domain/entities"
	"time"
)

type ValidationResult struct {
	IsValid          bool
	ValidationChecks map[string]CheckResult
	Errors           []string
	Warnings         []string
	ExpiresAt        *time.Time
}

type CheckResult struct {
	CheckName string
	Passed    bool
	Message   string
	Timestamp time.Time
}

type ExternalAccountValidationService interface {
	ValidateAccount(ctx context.Context, account *entities.ExternalAccount) (*ValidationResult, error)
	ShouldRevalidate(ctx context.Context, account *entities.ExternalAccount) bool
}

type validationService struct {
	// In a real implementation, you'd inject external validation providers
	// For IBAN validation, bank verification, fraud checks, etc.
}

func NewValidationService() ExternalAccountValidationService {
	return &validationService{}
}

func (s *validationService) ValidateAccount(ctx context.Context, account *entities.ExternalAccount) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:          true,
		ValidationChecks: make(map[string]CheckResult),
		Errors:           []string{},
		Warnings:         []string{},
	}

	now := time.Now()
	expiresAt := now.Add(90 * 24 * time.Hour) // Valid for 90 days
	result.ExpiresAt = &expiresAt

	// IBAN validation
	if account.IBAN != "" {
		if s.validateIBAN(account.IBAN) {
			result.ValidationChecks["ibanValidation"] = CheckResult{
				CheckName: "IBAN Validation",
				Passed:    true,
				Message:   "Valid IBAN format",
				Timestamp: now,
			}
		} else {
			result.IsValid = false
			result.ValidationChecks["ibanValidation"] = CheckResult{
				CheckName: "IBAN Validation",
				Passed:    false,
				Message:   "Invalid IBAN format",
				Timestamp: now,
			}
			result.Errors = append(result.Errors, "Invalid IBAN format")
		}
	}

	// Account number validation
	if account.AccountNumber != "" {
		if s.validateAccountNumber(account.AccountNumber) {
			result.ValidationChecks["accountNumberValidation"] = CheckResult{
				CheckName: "Account Number Validation",
				Passed:    true,
				Message:   "Valid account number format",
				Timestamp: now,
			}
		} else {
			result.IsValid = false
			result.ValidationChecks["accountNumberValidation"] = CheckResult{
				CheckName: "Account Number Validation",
				Passed:    false,
				Message:   "Invalid account number format",
				Timestamp: now,
			}
			result.Errors = append(result.Errors, "Invalid account number format")
		}
	}

	// Routing number validation (for US accounts)
	if account.RoutingNumber != "" {
		if s.validateRoutingNumber(account.RoutingNumber) {
			result.ValidationChecks["routingNumberValidation"] = CheckResult{
				CheckName: "Routing Number Validation",
				Passed:    true,
				Message:   "Valid routing number format",
				Timestamp: now,
			}
		} else {
			result.IsValid = false
			result.ValidationChecks["routingNumberValidation"] = CheckResult{
				CheckName: "Routing Number Validation",
				Passed:    false,
				Message:   "Invalid routing number format",
				Timestamp: now,
			}
			result.Errors = append(result.Errors, "Invalid routing number format")
		}
	}

	// TODO: Integrate with bank verification API
	// TODO: Integrate with fraud detection service
	// TODO: Check against sanctions lists

	// Basic currency validation
	if account.Currency == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Currency is required")
	}

	return result, nil
}

func (s *validationService) ShouldRevalidate(ctx context.Context, account *entities.ExternalAccount) bool {
	return account.NeedsRevalidation()
}

// Basic IBAN validation (simplified - real implementation should use proper IBAN algorithm)
func (s *validationService) validateIBAN(iban string) bool {
	// Basic format check: 2 letters (country code) + 2 digits (check digits) + up to 30 alphanumeric characters
	if len(iban) < 4 || len(iban) > 34 {
		return false
	}
	// TODO: Implement proper IBAN validation algorithm
	return true
}

// Basic account number validation
func (s *validationService) validateAccountNumber(accountNumber string) bool {
	// Basic check: not empty and reasonable length
	return len(accountNumber) > 0 && len(accountNumber) <= 255
}

// Basic routing number validation (US: 9 digits)
func (s *validationService) validateRoutingNumber(routingNumber string) bool {
	// US routing numbers are 9 digits
	if len(routingNumber) != 9 {
		return false
	}
	// TODO: Implement proper routing number checksum validation
	return true
}
