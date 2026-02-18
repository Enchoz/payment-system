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

type paymentRepository struct {
	db *DB
}

func NewPaymentRepository(db *DB) repositories.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	query := `
		INSERT INTO payments (
			id, paymentReference, requestId, fromUserId, toUserId, externalAccountId,
			amount, currency, exchangeRate, convertedAmount, convertedCurrency,
			status, paymentType, externalAccountNumber, externalBankName,
			failureReason, processedAt, completedAt, metadata,
			createdAt, updatedAt
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21
		)
	`

	metadataJSON, _ := mapToJSONB(payment.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		payment.PaymentReference,
		payment.RequestID,
		payment.FromUserID,
		payment.ToUserID,
		payment.ExternalAccountID,
		payment.Amount,
		payment.Currency,
		payment.ExchangeRate,
		payment.ConvertedAmount,
		payment.ConvertedCurrency,
		string(payment.Status),
		string(payment.PaymentType),
		payment.ExternalAccountNumber,
		payment.ExternalBankName,
		payment.FailureReason,
		payment.ProcessedAt,
		payment.CompletedAt,
		metadataJSON,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	query := `
		SELECT id, paymentReference, requestId, fromUserId, toUserId, externalAccountId,
			amount, currency, exchangeRate, convertedAmount, convertedCurrency,
			status, paymentType, externalAccountNumber, externalBankName,
			failureReason, processedAt, completedAt, metadata,
			createdAt, updatedAt
		FROM payments
		WHERE id = $1
	`

	payment := &entities.Payment{}
	var amountStr string
	var metadataJSON []byte
	var statusStr, paymentTypeStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.PaymentReference,
		&payment.RequestID,
		&payment.FromUserID,
		&payment.ToUserID,
		&payment.ExternalAccountID,
		&amountStr,
		&payment.Currency,
		&payment.ExchangeRate,
		&payment.ConvertedAmount,
		&payment.ConvertedCurrency,
		&statusStr,
		&paymentTypeStr,
		&payment.ExternalAccountNumber,
		&payment.ExternalBankName,
		&payment.FailureReason,
		&payment.ProcessedAt,
		&payment.CompletedAt,
		&metadataJSON,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	payment.Amount, _ = decimal.NewFromString(amountStr)
	payment.Status = entities.PaymentStatus(statusStr)
	payment.PaymentType = entities.PaymentType(paymentTypeStr)
	payment.Metadata, _ = jsonbToMap(metadataJSON)

	return payment, nil
}

func (r *paymentRepository) GetByPaymentReference(ctx context.Context, reference string) (*entities.Payment, error) {
	query := `
		SELECT id, paymentReference, requestId, fromUserId, toUserId, externalAccountId,
			amount, currency, exchangeRate, convertedAmount, convertedCurrency,
			status, paymentType, externalAccountNumber, externalBankName,
			failureReason, processedAt, completedAt, metadata,
			createdAt, updatedAt
		FROM payments
		WHERE paymentReference = $1
	`

	payment := &entities.Payment{}
	var amountStr string
	var metadataJSON []byte
	var statusStr, paymentTypeStr string

	err := r.db.QueryRowContext(ctx, query, reference).Scan(
		&payment.ID,
		&payment.PaymentReference,
		&payment.RequestID,
		&payment.FromUserID,
		&payment.ToUserID,
		&payment.ExternalAccountID,
		&amountStr,
		&payment.Currency,
		&payment.ExchangeRate,
		&payment.ConvertedAmount,
		&payment.ConvertedCurrency,
		&statusStr,
		&paymentTypeStr,
		&payment.ExternalAccountNumber,
		&payment.ExternalBankName,
		&payment.FailureReason,
		&payment.ProcessedAt,
		&payment.CompletedAt,
		&metadataJSON,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	payment.Amount, _ = decimal.NewFromString(amountStr)
	payment.Status = entities.PaymentStatus(statusStr)
	payment.PaymentType = entities.PaymentType(paymentTypeStr)
	payment.Metadata, _ = jsonbToMap(metadataJSON)

	return payment, nil
}

func (r *paymentRepository) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*entities.Payment, error) {
	query := `
		SELECT id, paymentReference, requestId, fromUserId, toUserId, externalAccountId,
			amount, currency, exchangeRate, convertedAmount, convertedCurrency,
			status, paymentType, externalAccountNumber, externalBankName,
			failureReason, processedAt, completedAt, metadata,
			createdAt, updatedAt
		FROM payments
		WHERE requestId = $1
	`

	payment := &entities.Payment{}
	var amountStr string
	var metadataJSON []byte
	var statusStr, paymentTypeStr string

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&payment.ID,
		&payment.PaymentReference,
		&payment.RequestID,
		&payment.FromUserID,
		&payment.ToUserID,
		&payment.ExternalAccountID,
		&amountStr,
		&payment.Currency,
		&payment.ExchangeRate,
		&payment.ConvertedAmount,
		&payment.ConvertedCurrency,
		&statusStr,
		&paymentTypeStr,
		&payment.ExternalAccountNumber,
		&payment.ExternalBankName,
		&payment.FailureReason,
		&payment.ProcessedAt,
		&payment.CompletedAt,
		&metadataJSON,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	payment.Amount, _ = decimal.NewFromString(amountStr)
	payment.Status = entities.PaymentStatus(statusStr)
	payment.PaymentType = entities.PaymentType(paymentTypeStr)
	payment.Metadata, _ = jsonbToMap(metadataJSON)

	return payment, nil
}

func (r *paymentRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Payment, error) {
	query := `
		SELECT id, paymentReference, requestId, fromUserId, toUserId, externalAccountId,
			amount, currency, exchangeRate, convertedAmount, convertedCurrency,
			status, paymentType, externalAccountNumber, externalBankName,
			failureReason, processedAt, completedAt, metadata,
			createdAt, updatedAt
		FROM payments
		WHERE fromUserId = $1 OR toUserId = $1
		ORDER BY createdAt DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}
	defer rows.Close()

	var payments []*entities.Payment
	for rows.Next() {
		payment := &entities.Payment{}
		var amountStr string
		var metadataJSON []byte
		var statusStr, paymentTypeStr string

		if err := rows.Scan(
			&payment.ID,
			&payment.PaymentReference,
			&payment.RequestID,
			&payment.FromUserID,
			&payment.ToUserID,
			&payment.ExternalAccountID,
			&amountStr,
			&payment.Currency,
			&payment.ExchangeRate,
			&payment.ConvertedAmount,
			&payment.ConvertedCurrency,
			&statusStr,
			&paymentTypeStr,
			&payment.ExternalAccountNumber,
			&payment.ExternalBankName,
			&payment.FailureReason,
			&payment.ProcessedAt,
			&payment.CompletedAt,
			&metadataJSON,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		payment.Amount, _ = decimal.NewFromString(amountStr)
		payment.Status = entities.PaymentStatus(statusStr)
		payment.PaymentType = entities.PaymentType(paymentTypeStr)
		payment.Metadata, _ = jsonbToMap(metadataJSON)

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	query := `
		UPDATE payments
		SET status = $2, failureReason = $3, completedAt = $4, metadata = $5, updatedAt = $6
		WHERE id = $1
	`

	metadataJSON, _ := mapToJSONB(payment.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		string(payment.Status),
		payment.FailureReason,
		payment.CompletedAt,
		metadataJSON,
		payment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	return nil
}
