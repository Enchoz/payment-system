package services

import (
	"context"
	"fmt"
	"payment-system/internal/domain/entities"
	"payment-system/internal/domain/repositories"
	"payment-system/internal/infrastructure/postgres"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentService interface {
	ProcessInternalPayment(ctx context.Context, req InternalPaymentRequest) (*entities.Payment, error)
	CreateExternalPaymentRequest(ctx context.Context, req ExternalPaymentRequest) (*entities.PaymentRequest, error)
	ProcessPaymentRequest(ctx context.Context, requestID uuid.UUID) (*entities.Payment, error)
}

type InternalPaymentRequest struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	Amount     decimal.Decimal
	Currency   string
	Reference  string
}

type ExternalPaymentRequest struct {
	FromUserID            uuid.UUID
	ExternalAccountID     *uuid.UUID
	Amount                decimal.Decimal
	Currency              string
	TargetCurrency        string
	ExternalAccountNumber string
	ExternalRoutingNumber string
	ExternalIBAN          string
	ExternalSwiftCode     string
	ExternalBankName       string
	ExternalAccountHolderName string
	ExternalCountryCode    string
	Reference              string
	IdempotencyKey         string
	SaveAccount            bool
}

type paymentService struct {
	db                    *postgres.DB
	userRepo              repositories.UserRepository
	accountRepo           repositories.AccountRepository
	externalAccountRepo   repositories.ExternalAccountRepository
	paymentRequestRepo    repositories.PaymentRequestRepository
	paymentRepo           repositories.PaymentRepository
	ledgerRepo            repositories.LedgerRepository
	paymentHoldRepo       repositories.PaymentHoldRepository
	exchangeRateService   ExchangeRateService
	validationService     ExternalAccountValidationService
}

func NewPaymentService(
	db *postgres.DB,
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	externalAccountRepo repositories.ExternalAccountRepository,
	paymentRequestRepo repositories.PaymentRequestRepository,
	paymentRepo repositories.PaymentRepository,
	ledgerRepo repositories.LedgerRepository,
	paymentHoldRepo repositories.PaymentHoldRepository,
	exchangeRateService ExchangeRateService,
	validationService ExternalAccountValidationService,
) PaymentService {
	return &paymentService{
		db:                  db,
		userRepo:            userRepo,
		accountRepo:         accountRepo,
		externalAccountRepo: externalAccountRepo,
		paymentRequestRepo:  paymentRequestRepo,
		paymentRepo:         paymentRepo,
		ledgerRepo:          ledgerRepo,
		paymentHoldRepo:     paymentHoldRepo,
		exchangeRateService: exchangeRateService,
		validationService:   validationService,
	}
}

func (s *paymentService) ProcessInternalPayment(ctx context.Context, req InternalPaymentRequest) (*entities.Payment, error) {
	// Validate users exist
	fromUser, err := s.userRepo.GetByID(ctx, req.FromUserID)
	if err != nil {
		return nil, fmt.Errorf("from user not found: %w", err)
	}

	toUser, err := s.userRepo.GetByID(ctx, req.ToUserID)
	if err != nil {
		return nil, fmt.Errorf("to user not found: %w", err)
	}

	// Get accounts
	fromAccount, err := s.accountRepo.GetByUserIDAndCurrency(ctx, fromUser.ID, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("from account not found: %w", err)
	}

	toAccount, err := s.accountRepo.GetByUserIDAndCurrency(ctx, toUser.ID, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("to account not found: %w", err)
	}

	// Check balance
	if !fromAccount.HasSufficientBalance(req.Amount) {
		return nil, entities.ErrInsufficientBalance
	}

	// Start transaction
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Debit from account
	if err := fromAccount.Debit(req.Amount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Update(ctx, fromAccount); err != nil {
		return nil, fmt.Errorf("failed to update from account: %w", err)
	}

	// Credit to account
	toAccount.Credit(req.Amount)
	if err := s.accountRepo.Update(ctx, toAccount); err != nil {
		return nil, fmt.Errorf("failed to update to account: %w", err)
	}

	// Create payment
	payment := entities.NewPayment(nil, req.FromUserID, req.Amount, req.Currency, entities.PaymentTypeInternal)
	payment.ToUserID = &req.ToUserID
	payment.PaymentReference = req.Reference
	payment.MarkAsCompleted()

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Create ledger entries
	fromEntry := entities.NewLedgerEntry(
		payment.ID,
		fromAccount.ID,
		entities.EntryTypeDebit,
		req.Amount,
		req.Currency,
		fromAccount.Balance,
	)
	if err := s.ledgerRepo.Create(ctx, fromEntry); err != nil {
		return nil, fmt.Errorf("failed to create debit ledger entry: %w", err)
	}

	toEntry := entities.NewLedgerEntry(
		payment.ID,
		toAccount.ID,
		entities.EntryTypeCredit,
		req.Amount,
		req.Currency,
		toAccount.Balance,
	)
	if err := s.ledgerRepo.Create(ctx, toEntry); err != nil {
		return nil, fmt.Errorf("failed to create credit ledger entry: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return payment, nil
}

func (s *paymentService) CreateExternalPaymentRequest(ctx context.Context, req ExternalPaymentRequest) (*entities.PaymentRequest, error) {
	// Check idempotency
	if req.IdempotencyKey != "" {
		existing, err := s.paymentRequestRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(ctx, req.FromUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Create payment request
	paymentRequest := entities.NewExternalPaymentRequest(req.FromUserID, req.Amount, req.Currency, req.TargetCurrency)
	paymentRequest.RequestReference = req.Reference
	paymentRequest.IdempotencyKey = req.IdempotencyKey

	// Handle external account
	if req.ExternalAccountID != nil {
		// Use saved account
		externalAccount, err := s.externalAccountRepo.GetByID(ctx, *req.ExternalAccountID)
		if err != nil {
			return nil, fmt.Errorf("external account not found: %w", err)
		}
		paymentRequest.ExternalAccountID = req.ExternalAccountID
		paymentRequest.ExternalAccountNumber = externalAccount.AccountNumber
		paymentRequest.ExternalIBAN = externalAccount.IBAN
		paymentRequest.ExternalRoutingNumber = externalAccount.RoutingNumber
		paymentRequest.ExternalBankName = externalAccount.BankName
		paymentRequest.ExternalAccountHolderName = externalAccount.AccountHolderName
		paymentRequest.ExternalCountryCode = externalAccount.CountryCode
	} else {
		// New account details provided
		paymentRequest.ExternalAccountNumber = req.ExternalAccountNumber
		paymentRequest.ExternalRoutingNumber = req.ExternalRoutingNumber
		paymentRequest.ExternalIBAN = req.ExternalIBAN
		paymentRequest.ExternalSwiftCode = req.ExternalSwiftCode
		paymentRequest.ExternalBankName = req.ExternalBankName
		paymentRequest.ExternalAccountHolderName = req.ExternalAccountHolderName
		paymentRequest.ExternalCountryCode = req.ExternalCountryCode

		// Optionally save the account
		if req.SaveAccount {
			externalAccount := entities.NewExternalAccount(req.FromUserID, req.ExternalAccountNumber, req.TargetCurrency)
			externalAccount.RoutingNumber = req.ExternalRoutingNumber
			externalAccount.IBAN = req.ExternalIBAN
			externalAccount.SwiftCode = req.ExternalSwiftCode
			externalAccount.BankName = req.ExternalBankName
			externalAccount.AccountHolderName = req.ExternalAccountHolderName
			externalAccount.CountryCode = req.ExternalCountryCode

			if err := s.externalAccountRepo.Create(ctx, externalAccount); err != nil {
				return nil, fmt.Errorf("failed to save external account: %w", err)
			}
			paymentRequest.ExternalAccountID = &externalAccount.ID
		}
	}

	// Save payment request
	if err := s.paymentRequestRepo.Create(ctx, paymentRequest); err != nil {
		return nil, fmt.Errorf("failed to create payment request: %w", err)
	}

	// Start validation process (async in production)
	paymentRequest.MarkAsValidating()
	if err := s.paymentRequestRepo.Update(ctx, paymentRequest); err != nil {
		return nil, fmt.Errorf("failed to update payment request status: %w", err)
	}

	// Validate external account
	var externalAccount *entities.ExternalAccount
	if paymentRequest.ExternalAccountID != nil {
		externalAccount, err = s.externalAccountRepo.GetByID(ctx, *paymentRequest.ExternalAccountID)
		if err != nil {
			return nil, fmt.Errorf("failed to get external account: %w", err)
		}
	} else {
		// Create temporary account entity for validation
		externalAccount = entities.NewExternalAccount(req.FromUserID, req.ExternalAccountNumber, req.TargetCurrency)
		externalAccount.RoutingNumber = req.ExternalRoutingNumber
		externalAccount.IBAN = req.ExternalIBAN
		externalAccount.BankName = req.ExternalBankName
		externalAccount.CountryCode = req.ExternalCountryCode
	}

	validationResult, err := s.validationService.ValidateAccount(ctx, externalAccount)
	if err != nil {
		paymentRequest.MarkAsFailed(fmt.Sprintf("validation error: %v", err))
		s.paymentRequestRepo.Update(ctx, paymentRequest)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Store validation results
	validationResultsMap := make(map[string]interface{})
	for k, v := range validationResult.ValidationChecks {
		validationResultsMap[k] = map[string]interface{}{
			"checkName": v.CheckName,
			"passed":    v.Passed,
			"message":   v.Message,
			"timestamp": v.Timestamp,
		}
	}
	paymentRequest.SetValidationResults(validationResultsMap)

	if !validationResult.IsValid {
		validationErrorsMap := make(map[string]interface{})
		validationErrorsMap["errors"] = validationResult.Errors
		paymentRequest.SetValidationErrors(validationErrorsMap)
		paymentRequest.MarkAsFailed("Account validation failed")
		s.paymentRequestRepo.Update(ctx, paymentRequest)
		return nil, fmt.Errorf("account validation failed: %v", validationResult.Errors)
	}

	// Update external account validation status if saved
	if paymentRequest.ExternalAccountID != nil {
		externalAccount.MarkAsValid(validationResult.ExpiresAt)
		if err := s.externalAccountRepo.Update(ctx, externalAccount); err != nil {
			return nil, fmt.Errorf("failed to update external account validation: %w", err)
		}
	}

	// Handle currency conversion if needed
	if req.Currency != req.TargetCurrency {
		convertedAmount, err := s.exchangeRateService.ConvertAmount(ctx, req.Amount, req.Currency, req.TargetCurrency)
		if err != nil {
			paymentRequest.MarkAsFailed(fmt.Sprintf("exchange rate error: %v", err))
			s.paymentRequestRepo.Update(ctx, paymentRequest)
			return nil, fmt.Errorf("failed to convert currency: %w", err)
		}

		rate, err := s.exchangeRateService.GetExchangeRate(ctx, req.Currency, req.TargetCurrency)
		if err != nil {
			paymentRequest.MarkAsFailed(fmt.Sprintf("exchange rate error: %v", err))
			s.paymentRequestRepo.Update(ctx, paymentRequest)
			return nil, fmt.Errorf("failed to get exchange rate: %w", err)
		}

		paymentRequest.SetExchangeRate(rate.Rate, *convertedAmount)
	}

	paymentRequest.MarkAsValidated()
	if err := s.paymentRequestRepo.Update(ctx, paymentRequest); err != nil {
		return nil, fmt.Errorf("failed to update payment request: %w", err)
	}

	return paymentRequest, nil
}

func (s *paymentService) ProcessPaymentRequest(ctx context.Context, requestID uuid.UUID) (*entities.Payment, error) {
	// Get payment request
	request, err := s.paymentRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("payment request not found: %w", err)
	}

	if request.Status != entities.PaymentRequestStatusValidated {
		return nil, fmt.Errorf("payment request is not validated, current status: %s", request.Status)
	}

	// Get from account
	fromAccount, err := s.accountRepo.GetByUserIDAndCurrency(ctx, request.FromUserID, request.Currency)
	if err != nil {
		return nil, fmt.Errorf("from account not found: %w", err)
	}

	// Determine amount to debit
	debitAmount := request.Amount
	if request.ConvertedAmount != nil {
		// For external payments, we debit the original currency
		debitAmount = request.Amount
	}

	// Check balance
	if !fromAccount.HasSufficientBalance(debitAmount) {
		request.MarkAsFailed("insufficient balance")
		s.paymentRequestRepo.Update(ctx, request)
		return nil, entities.ErrInsufficientBalance
	}

	// Start transaction
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Place hold
	hold := entities.NewPaymentHold(fromAccount.ID, debitAmount, request.Currency, &request.ID)
	if err := fromAccount.PlaceHold(debitAmount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Update(ctx, fromAccount); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}
	if err := s.paymentHoldRepo.Create(ctx, hold); err != nil {
		return nil, fmt.Errorf("failed to create payment hold: %w", err)
	}

	request.MarkAsProcessing()
	if err := s.paymentRequestRepo.Update(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to update payment request: %w", err)
	}

	// Create payment
	payment := entities.NewPayment(&request.ID, request.FromUserID, request.Amount, request.Currency, request.PaymentType)
	payment.PaymentReference = request.RequestReference
	if request.ExchangeRate != nil && request.ConvertedAmount != nil {
		payment.SetExchangeDetails(*request.ExchangeRate, *request.ConvertedAmount, request.TargetCurrency)
	}
	if request.ExternalAccountID != nil {
		payment.ExternalAccountID = request.ExternalAccountID
	}
	payment.ExternalAccountNumber = request.ExternalAccountNumber
	payment.ExternalBankName = request.ExternalBankName

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Execute payment
	if err := fromAccount.Debit(debitAmount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Update(ctx, fromAccount); err != nil {
		return nil, fmt.Errorf("failed to update from account: %w", err)
	}

	// Create ledger entry for debit
	debitEntry := entities.NewLedgerEntry(
		payment.ID,
		fromAccount.ID,
		entities.EntryTypeDebit,
		debitAmount,
		request.Currency,
		fromAccount.Balance,
	)
	if err := s.ledgerRepo.Create(ctx, debitEntry); err != nil {
		return nil, fmt.Errorf("failed to create debit ledger entry: %w", err)
	}

	// For external payments, credit a system liability account
	// For internal payments, credit the recipient account
	if request.PaymentType == entities.PaymentTypeInternal && request.ToUserID != nil {
		toAccount, err := s.accountRepo.GetByUserIDAndCurrency(ctx, *request.ToUserID, request.Currency)
		if err != nil {
			return nil, fmt.Errorf("to account not found: %w", err)
		}
		toAccount.Credit(request.Amount)
		if err := s.accountRepo.Update(ctx, toAccount); err != nil {
			return nil, fmt.Errorf("failed to update to account: %w", err)
		}

		creditEntry := entities.NewLedgerEntry(
			payment.ID,
			toAccount.ID,
			entities.EntryTypeCredit,
			request.Amount,
			request.Currency,
			toAccount.Balance,
		)
		if err := s.ledgerRepo.Create(ctx, creditEntry); err != nil {
			return nil, fmt.Errorf("failed to create credit ledger entry: %w", err)
		}
		payment.ToUserID = request.ToUserID
	}

	// Release hold
	hold.Release()
	if err := s.paymentHoldRepo.Update(ctx, hold); err != nil {
		return nil, fmt.Errorf("failed to release payment hold: %w", err)
	}

	payment.MarkAsCompleted()
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	request.MarkAsCompleted()
	if err := s.paymentRequestRepo.Update(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to update payment request: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return payment, nil
}
