package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"payment-system/internal/domain/entities"
	"payment-system/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type paymentRequestRepository struct {
	db *DB
}

func NewPaymentRequestRepository(db *DB) repositories.PaymentRequestRepository {
	return &paymentRequestRepository{db: db}
}

func (r *paymentRequestRepository) Create(ctx context.Context, request *entities.PaymentRequest) error {
	query := `
		INSERT INTO payment_requests (
			id, requestReference, fromUserId, externalAccountId, toUserId,
			amount, currency, targetCurrency,
			externalAccountNumber, externalRoutingNumber, externalIban, externalSwiftCode,
			externalBankName, externalAccountHolderName, externalCountryCode,
			status, paymentType, validationResults, validationErrors,
			requiresManualReview, exchangeRate, convertedAmount,
			failureReason, processedAt, metadata, idempotencyKey,
			createdAt, updatedAt
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28
		)
	`

	validationResultsJSON, _ := mapToJSONB(request.ValidationResults)
	validationErrorsJSON, _ := mapToJSONB(request.ValidationErrors)
	metadataJSON, _ := mapToJSONB(request.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		request.ID,
		request.RequestReference,
		request.FromUserID,
		request.ExternalAccountID,
		request.ToUserID,
		request.Amount,
		request.Currency,
		request.TargetCurrency,
		request.ExternalAccountNumber,
		request.ExternalRoutingNumber,
		request.ExternalIBAN,
		request.ExternalSwiftCode,
		request.ExternalBankName,
		request.ExternalAccountHolderName,
		request.ExternalCountryCode,
		string(request.Status),
		string(request.PaymentType),
		validationResultsJSON,
		validationErrorsJSON,
		request.RequiresManualReview,
		request.ExchangeRate,
		request.ConvertedAmount,
		request.FailureReason,
		request.ProcessedAt,
		metadataJSON,
		request.IdempotencyKey,
		request.CreatedAt,
		request.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment request: %w", err)
	}
	return nil
}

func (r *paymentRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.PaymentRequest, error) {
	query := `
		SELECT id, requestReference, fromUserId, externalAccountId, toUserId,
			amount, currency, targetCurrency,
			externalAccountNumber, externalRoutingNumber, externalIban, externalSwiftCode,
			externalBankName, externalAccountHolderName, externalCountryCode,
			status, paymentType, validationResults, validationErrors,
			requiresManualReview, exchangeRate, convertedAmount,
			failureReason, processedAt, metadata, idempotencyKey,
			createdAt, updatedAt
		FROM payment_requests
		WHERE id = $1
	`

	request := &entities.PaymentRequest{}
	var amountStr string
	var validationResultsJSON, validationErrorsJSON, metadataJSON []byte
	var statusStr, paymentTypeStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.RequestReference,
		&request.FromUserID,
		&request.ExternalAccountID,
		&request.ToUserID,
		&amountStr,
		&request.Currency,
		&request.TargetCurrency,
		&request.ExternalAccountNumber,
		&request.ExternalRoutingNumber,
		&request.ExternalIBAN,
		&request.ExternalSwiftCode,
		&request.ExternalBankName,
		&request.ExternalAccountHolderName,
		&request.ExternalCountryCode,
		&statusStr,
		&paymentTypeStr,
		&validationResultsJSON,
		&validationErrorsJSON,
		&request.RequiresManualReview,
		&request.ExchangeRate,
		&request.ConvertedAmount,
		&request.FailureReason,
		&request.ProcessedAt,
		&metadataJSON,
		&request.IdempotencyKey,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment request not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment request: %w", err)
	}

	request.Amount, _ = decimal.NewFromString(amountStr)
	request.Status = entities.PaymentRequestStatus(statusStr)
	request.PaymentType = entities.PaymentType(paymentTypeStr)
	request.ValidationResults, _ = jsonbToMap(validationResultsJSON)
	request.ValidationErrors, _ = jsonbToMap(validationErrorsJSON)
	request.Metadata, _ = jsonbToMap(metadataJSON)

	return request, nil
}

func (r *paymentRequestRepository) GetByRequestReference(ctx context.Context, reference string) (*entities.PaymentRequest, error) {
	query := `
		SELECT id, requestReference, fromUserId, externalAccountId, toUserId,
			amount, currency, targetCurrency,
			externalAccountNumber, externalRoutingNumber, externalIban, externalSwiftCode,
			externalBankName, externalAccountHolderName, externalCountryCode,
			status, paymentType, validationResults, validationErrors,
			requiresManualReview, exchangeRate, convertedAmount,
			failureReason, processedAt, metadata, idempotencyKey,
			createdAt, updatedAt
		FROM payment_requests
		WHERE requestReference = $1
	`

	request := &entities.PaymentRequest{}
	var amountStr string
	var validationResultsJSON, validationErrorsJSON, metadataJSON []byte
	var statusStr, paymentTypeStr string

	err := r.db.QueryRowContext(ctx, query, reference).Scan(
		&request.ID,
		&request.RequestReference,
		&request.FromUserID,
		&request.ExternalAccountID,
		&request.ToUserID,
		&amountStr,
		&request.Currency,
		&request.TargetCurrency,
		&request.ExternalAccountNumber,
		&request.ExternalRoutingNumber,
		&request.ExternalIBAN,
		&request.ExternalSwiftCode,
		&request.ExternalBankName,
		&request.ExternalAccountHolderName,
		&request.ExternalCountryCode,
		&statusStr,
		&paymentTypeStr,
		&validationResultsJSON,
		&validationErrorsJSON,
		&request.RequiresManualReview,
		&request.ExchangeRate,
		&request.ConvertedAmount,
		&request.FailureReason,
		&request.ProcessedAt,
		&metadataJSON,
		&request.IdempotencyKey,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment request not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment request: %w", err)
	}

	request.Amount, _ = decimal.NewFromString(amountStr)
	request.Status = entities.PaymentRequestStatus(statusStr)
	request.PaymentType = entities.PaymentType(paymentTypeStr)
	request.ValidationResults, _ = jsonbToMap(validationResultsJSON)
	request.ValidationErrors, _ = jsonbToMap(validationErrorsJSON)
	request.Metadata, _ = jsonbToMap(metadataJSON)

	return request, nil
}

func (r *paymentRequestRepository) GetByIdempotencyKey(ctx context.Context, key string) (*entities.PaymentRequest, error) {
	query := `
		SELECT id, requestReference, fromUserId, externalAccountId, toUserId,
			amount, currency, targetCurrency,
			externalAccountNumber, externalRoutingNumber, externalIban, externalSwiftCode,
			externalBankName, externalAccountHolderName, externalCountryCode,
			status, paymentType, validationResults, validationErrors,
			requiresManualReview, exchangeRate, convertedAmount,
			failureReason, processedAt, metadata, idempotencyKey,
			createdAt, updatedAt
		FROM payment_requests
		WHERE idempotencyKey = $1
	`

	request := &entities.PaymentRequest{}
	var amountStr string
	var validationResultsJSON, validationErrorsJSON, metadataJSON []byte
	var statusStr, paymentTypeStr string

	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&request.ID,
		&request.RequestReference,
		&request.FromUserID,
		&request.ExternalAccountID,
		&request.ToUserID,
		&amountStr,
		&request.Currency,
		&request.TargetCurrency,
		&request.ExternalAccountNumber,
		&request.ExternalRoutingNumber,
		&request.ExternalIBAN,
		&request.ExternalSwiftCode,
		&request.ExternalBankName,
		&request.ExternalAccountHolderName,
		&request.ExternalCountryCode,
		&statusStr,
		&paymentTypeStr,
		&validationResultsJSON,
		&validationErrorsJSON,
		&request.RequiresManualReview,
		&request.ExchangeRate,
		&request.ConvertedAmount,
		&request.FailureReason,
		&request.ProcessedAt,
		&metadataJSON,
		&request.IdempotencyKey,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment request not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment request: %w", err)
	}

	request.Amount, _ = decimal.NewFromString(amountStr)
	request.Status = entities.PaymentRequestStatus(statusStr)
	request.PaymentType = entities.PaymentType(paymentTypeStr)
	request.ValidationResults, _ = jsonbToMap(validationResultsJSON)
	request.ValidationErrors, _ = jsonbToMap(validationErrorsJSON)
	request.Metadata, _ = jsonbToMap(metadataJSON)

	return request, nil
}

func (r *paymentRequestRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.PaymentRequest, error) {
	query := `
		SELECT id, requestReference, fromUserId, externalAccountId, toUserId,
			amount, currency, targetCurrency,
			externalAccountNumber, externalRoutingNumber, externalIban, externalSwiftCode,
			externalBankName, externalAccountHolderName, externalCountryCode,
			status, paymentType, validationResults, validationErrors,
			requiresManualReview, exchangeRate, convertedAmount,
			failureReason, processedAt, metadata, idempotencyKey,
			createdAt, updatedAt
		FROM payment_requests
		WHERE fromUserId = $1 OR toUserId = $1
		ORDER BY createdAt DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment requests: %w", err)
	}
	defer rows.Close()

	var requests []*entities.PaymentRequest
	for rows.Next() {
		request := &entities.PaymentRequest{}
		var amountStr string
		var validationResultsJSON, validationErrorsJSON, metadataJSON []byte
		var statusStr, paymentTypeStr string

		if err := rows.Scan(
			&request.ID,
			&request.RequestReference,
			&request.FromUserID,
			&request.ExternalAccountID,
			&request.ToUserID,
			&amountStr,
			&request.Currency,
			&request.TargetCurrency,
			&request.ExternalAccountNumber,
			&request.ExternalRoutingNumber,
			&request.ExternalIBAN,
			&request.ExternalSwiftCode,
			&request.ExternalBankName,
			&request.ExternalAccountHolderName,
			&request.ExternalCountryCode,
			&statusStr,
			&paymentTypeStr,
			&validationResultsJSON,
			&validationErrorsJSON,
			&request.RequiresManualReview,
			&request.ExchangeRate,
			&request.ConvertedAmount,
			&request.FailureReason,
			&request.ProcessedAt,
			&metadataJSON,
			&request.IdempotencyKey,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment request: %w", err)
		}

		request.Amount, _ = decimal.NewFromString(amountStr)
		request.Status = entities.PaymentRequestStatus(statusStr)
		request.PaymentType = entities.PaymentType(paymentTypeStr)
		request.ValidationResults, _ = jsonbToMap(validationResultsJSON)
		request.ValidationErrors, _ = jsonbToMap(validationErrorsJSON)
		request.Metadata, _ = jsonbToMap(metadataJSON)

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *paymentRequestRepository) Update(ctx context.Context, request *entities.PaymentRequest) error {
	query := `
		UPDATE payment_requests
		SET status = $2, validationResults = $3, validationErrors = $4,
			requiresManualReview = $5, exchangeRate = $6, convertedAmount = $7,
			failureReason = $8, processedAt = $9, metadata = $10, updatedAt = $11
		WHERE id = $1
	`

	validationResultsJSON, _ := mapToJSONB(request.ValidationResults)
	validationErrorsJSON, _ := mapToJSONB(request.ValidationErrors)
	metadataJSON, _ := mapToJSONB(request.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		request.ID,
		string(request.Status),
		validationResultsJSON,
		validationErrorsJSON,
		request.RequiresManualReview,
		request.ExchangeRate,
		request.ConvertedAmount,
		request.FailureReason,
		request.ProcessedAt,
		metadataJSON,
		request.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update payment request: %w", err)
	}
	return nil
}

func (r *paymentRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.PaymentRequestStatus) error {
	query := `UPDATE payment_requests SET status = $2, updatedAt = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, string(status))
	if err != nil {
		return fmt.Errorf("failed to update payment request status: %w", err)
	}
	return nil
}
