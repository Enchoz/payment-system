package postgres

import (
	"context"
	"fmt"
	"payment-system/internal/domain/entities"
	"payment-system/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ledgerRepository struct {
	db *DB
}

func NewLedgerRepository(db *DB) repositories.LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) Create(ctx context.Context, entry *entities.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (
			id, paymentId, accountId, entryType, amount, currency, balanceAfter, createdAt
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID,
		entry.PaymentID,
		entry.AccountID,
		string(entry.EntryType),
		entry.Amount,
		entry.Currency,
		entry.BalanceAfter,
		entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create ledger entry: %w", err)
	}
	return nil
}

func (r *ledgerRepository) GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*entities.LedgerEntry, error) {
	query := `
		SELECT id, paymentId, accountId, entryType, amount, currency, balanceAfter, createdAt
		FROM ledger_entries
		WHERE paymentId = $1
		ORDER BY createdAt
	`

	rows, err := r.db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entries: %w", err)
	}
	defer rows.Close()

	var entries []*entities.LedgerEntry
	for rows.Next() {
		entry := &entities.LedgerEntry{}
		var amountStr, balanceAfterStr string
		var entryTypeStr string

		if err := rows.Scan(
			&entry.ID,
			&entry.PaymentID,
			&entry.AccountID,
			&entryTypeStr,
			&amountStr,
			&entry.Currency,
			&balanceAfterStr,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ledger entry: %w", err)
		}

		entry.EntryType = entities.EntryType(entryTypeStr)
		entry.Amount, _ = decimal.NewFromString(amountStr)
		entry.BalanceAfter, _ = decimal.NewFromString(balanceAfterStr)

		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *ledgerRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entities.LedgerEntry, error) {
	query := `
		SELECT id, paymentId, accountId, entryType, amount, currency, balanceAfter, createdAt
		FROM ledger_entries
		WHERE accountId = $1
		ORDER BY createdAt DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entries: %w", err)
	}
	defer rows.Close()

	var entries []*entities.LedgerEntry
	for rows.Next() {
		entry := &entities.LedgerEntry{}
		var amountStr, balanceAfterStr string
		var entryTypeStr string

		if err := rows.Scan(
			&entry.ID,
			&entry.PaymentID,
			&entry.AccountID,
			&entryTypeStr,
			&amountStr,
			&entry.Currency,
			&balanceAfterStr,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ledger entry: %w", err)
		}

		entry.EntryType = entities.EntryType(entryTypeStr)
		entry.Amount, _ = decimal.NewFromString(amountStr)
		entry.BalanceAfter, _ = decimal.NewFromString(balanceAfterStr)

		entries = append(entries, entry)
	}

	return entries, nil
}
