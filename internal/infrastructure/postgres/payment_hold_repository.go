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

type paymentHoldRepository struct {
	db *DB
}

func NewPaymentHoldRepository(db *DB) repositories.PaymentHoldRepository {
	return &paymentHoldRepository{db: db}
}

func (r *paymentHoldRepository) Create(ctx context.Context, hold *entities.PaymentHold) error {
	query := `
		INSERT INTO payment_holds (
			id, accountId, requestId, paymentId, amount, currency, expiresAt, releasedAt, createdAt
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		hold.ID,
		hold.AccountID,
		hold.RequestID,
		hold.PaymentID,
		hold.Amount,
		hold.Currency,
		hold.ExpiresAt,
		hold.ReleasedAt,
		hold.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment hold: %w", err)
	}
	return nil
}

func (r *paymentHoldRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.PaymentHold, error) {
	query := `
		SELECT id, accountId, requestId, paymentId, amount, currency, expiresAt, releasedAt, createdAt
		FROM payment_holds
		WHERE id = $1
	`

	hold := &entities.PaymentHold{}
	var amountStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hold.ID,
		&hold.AccountID,
		&hold.RequestID,
		&hold.PaymentID,
		&amountStr,
		&hold.Currency,
		&hold.ExpiresAt,
		&hold.ReleasedAt,
		&hold.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment hold not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment hold: %w", err)
	}

	hold.Amount, _ = decimal.NewFromString(amountStr)
	return hold, nil
}

func (r *paymentHoldRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entities.PaymentHold, error) {
	query := `
		SELECT id, accountId, requestId, paymentId, amount, currency, expiresAt, releasedAt, createdAt
		FROM payment_holds
		WHERE accountId = $1 AND releasedAt IS NULL
		ORDER BY createdAt DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment holds: %w", err)
	}
	defer rows.Close()

	var holds []*entities.PaymentHold
	for rows.Next() {
		hold := &entities.PaymentHold{}
		var amountStr string

		if err := rows.Scan(
			&hold.ID,
			&hold.AccountID,
			&hold.RequestID,
			&hold.PaymentID,
			&amountStr,
			&hold.Currency,
			&hold.ExpiresAt,
			&hold.ReleasedAt,
			&hold.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment hold: %w", err)
		}

		hold.Amount, _ = decimal.NewFromString(amountStr)
		holds = append(holds, hold)
	}

	return holds, nil
}

func (r *paymentHoldRepository) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*entities.PaymentHold, error) {
	query := `
		SELECT id, accountId, requestId, paymentId, amount, currency, expiresAt, releasedAt, createdAt
		FROM payment_holds
		WHERE requestId = $1
	`

	hold := &entities.PaymentHold{}
	var amountStr string

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&hold.ID,
		&hold.AccountID,
		&hold.RequestID,
		&hold.PaymentID,
		&amountStr,
		&hold.Currency,
		&hold.ExpiresAt,
		&hold.ReleasedAt,
		&hold.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment hold not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment hold: %w", err)
	}

	hold.Amount, _ = decimal.NewFromString(amountStr)
	return hold, nil
}

func (r *paymentHoldRepository) GetByPaymentID(ctx context.Context, paymentID uuid.UUID) (*entities.PaymentHold, error) {
	query := `
		SELECT id, accountId, requestId, paymentId, amount, currency, expiresAt, releasedAt, createdAt
		FROM payment_holds
		WHERE paymentId = $1
	`

	hold := &entities.PaymentHold{}
	var amountStr string

	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(
		&hold.ID,
		&hold.AccountID,
		&hold.RequestID,
		&hold.PaymentID,
		&amountStr,
		&hold.Currency,
		&hold.ExpiresAt,
		&hold.ReleasedAt,
		&hold.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment hold not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment hold: %w", err)
	}

	hold.Amount, _ = decimal.NewFromString(amountStr)
	return hold, nil
}

func (r *paymentHoldRepository) Update(ctx context.Context, hold *entities.PaymentHold) error {
	query := `
		UPDATE payment_holds
		SET releasedAt = $2
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, hold.ID, hold.ReleasedAt)
	if err != nil {
		return fmt.Errorf("failed to update payment hold: %w", err)
	}
	return nil
}

func (r *paymentHoldRepository) Release(ctx context.Context, holdID uuid.UUID) error {
	query := `UPDATE payment_holds SET releasedAt = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, holdID)
	if err != nil {
		return fmt.Errorf("failed to release payment hold: %w", err)
	}
	return nil
}
