package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"payment-system/internal/domain/entities"
	"payment-system/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type accountRepository struct {
	db *DB
}

func NewAccountRepository(db *DB) repositories.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *entities.Account) error {
	query := `
		INSERT INTO accounts (id, userId, currency, balance, availableBalance, createdAt, updatedAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.Currency,
		account.Balance,
		account.AvailableBalance,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	query := `
		SELECT id, userId, currency, balance, availableBalance, createdAt, updatedAt
		FROM accounts
		WHERE id = $1
	`
	account := &entities.Account{}
	var balanceStr, availableBalanceStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Currency,
		&balanceStr,
		&availableBalanceStr,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, entities.ErrAccountNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	account.Balance, _ = decimal.NewFromString(balanceStr)
	account.AvailableBalance, _ = decimal.NewFromString(availableBalanceStr)
	return account, nil
}

func (r *accountRepository) GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*entities.Account, error) {
	query := `
		SELECT id, userId, currency, balance, availableBalance, createdAt, updatedAt
		FROM accounts
		WHERE userId = $1 AND currency = $2
	`
	account := &entities.Account{}
	var balanceStr, availableBalanceStr string
	err := r.db.QueryRowContext(ctx, query, userID, currency).Scan(
		&account.ID,
		&account.UserID,
		&account.Currency,
		&balanceStr,
		&availableBalanceStr,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, entities.ErrAccountNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	account.Balance, _ = decimal.NewFromString(balanceStr)
	account.AvailableBalance, _ = decimal.NewFromString(availableBalanceStr)
	return account, nil
}

func (r *accountRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Account, error) {
	query := `
		SELECT id, userId, currency, balance, availableBalance, createdAt, updatedAt
		FROM accounts
		WHERE userId = $1
		ORDER BY currency
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*entities.Account
	for rows.Next() {
		account := &entities.Account{}
		var balanceStr, availableBalanceStr string
		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Currency,
			&balanceStr,
			&availableBalanceStr,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		account.Balance, _ = decimal.NewFromString(balanceStr)
		account.AvailableBalance, _ = decimal.NewFromString(availableBalanceStr)
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, account *entities.Account) error {
	query := `
		UPDATE accounts
		SET balance = $2, availableBalance = $3, updatedAt = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.Balance,
		account.AvailableBalance,
		account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, accountID uuid.UUID, balance, availableBalance interface{}) error {
	query := `
		UPDATE accounts
		SET balance = $2, availableBalance = $3, updatedAt = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, accountID, balance, availableBalance)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

// Helper function to convert JSONB to map
func jsonbToMap(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return make(map[string]interface{}), nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Helper function to convert map to JSONB
func mapToJSONB(m map[string]interface{}) ([]byte, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}
