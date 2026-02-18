package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"payment-system/internal/domain/entities"
	"payment-system/internal/domain/repositories"

	"github.com/google/uuid"
)

type externalAccountRepository struct {
	db *DB
}

func NewExternalAccountRepository(db *DB) repositories.ExternalAccountRepository {
	return &externalAccountRepository{db: db}
}

func (r *externalAccountRepository) Create(ctx context.Context, account *entities.ExternalAccount) error {
	query := `
		INSERT INTO external_accounts (
			id, userId, accountNumber, routingNumber, iban, swiftCode,
			bankName, accountHolderName, currency, accountType, countryCode,
			validationStatus, validatedAt, validationExpiresAt, lastValidatedAt,
			nickname, isDefault, metadata, createdAt, updatedAt
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
	`

	metadataJSON, _ := mapToJSONB(account.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.AccountNumber,
		account.RoutingNumber,
		account.IBAN,
		account.SwiftCode,
		account.BankName,
		account.AccountHolderName,
		account.Currency,
		account.AccountType,
		account.CountryCode,
		string(account.ValidationStatus),
		account.ValidatedAt,
		account.ValidationExpiresAt,
		account.LastValidatedAt,
		account.Nickname,
		account.IsDefault,
		metadataJSON,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create external account: %w", err)
	}
	return nil
}

func (r *externalAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ExternalAccount, error) {
	query := `
		SELECT id, userId, accountNumber, routingNumber, iban, swiftCode,
			bankName, accountHolderName, currency, accountType, countryCode,
			validationStatus, validatedAt, validationExpiresAt, lastValidatedAt,
			nickname, isDefault, metadata, createdAt, updatedAt
		FROM external_accounts
		WHERE id = $1
	`

	account := &entities.ExternalAccount{}
	var validationStatusStr string
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.RoutingNumber,
		&account.IBAN,
		&account.SwiftCode,
		&account.BankName,
		&account.AccountHolderName,
		&account.Currency,
		&account.AccountType,
		&account.CountryCode,
		&validationStatusStr,
		&account.ValidatedAt,
		&account.ValidationExpiresAt,
		&account.LastValidatedAt,
		&account.Nickname,
		&account.IsDefault,
		&metadataJSON,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("external account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get external account: %w", err)
	}

	account.ValidationStatus = entities.ValidationStatus(validationStatusStr)
	account.Metadata, _ = jsonbToMap(metadataJSON)

	return account, nil
}

func (r *externalAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ExternalAccount, error) {
	query := `
		SELECT id, userId, accountNumber, routingNumber, iban, swiftCode,
			bankName, accountHolderName, currency, accountType, countryCode,
			validationStatus, validatedAt, validationExpiresAt, lastValidatedAt,
			nickname, isDefault, metadata, createdAt, updatedAt
		FROM external_accounts
		WHERE userId = $1
		ORDER BY isDefault DESC, createdAt DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get external accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*entities.ExternalAccount
	for rows.Next() {
		account := &entities.ExternalAccount{}
		var validationStatusStr string
		var metadataJSON []byte

		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.AccountNumber,
			&account.RoutingNumber,
			&account.IBAN,
			&account.SwiftCode,
			&account.BankName,
			&account.AccountHolderName,
			&account.Currency,
			&account.AccountType,
			&account.CountryCode,
			&validationStatusStr,
			&account.ValidatedAt,
			&account.ValidationExpiresAt,
			&account.LastValidatedAt,
			&account.Nickname,
			&account.IsDefault,
			&metadataJSON,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan external account: %w", err)
		}

		account.ValidationStatus = entities.ValidationStatus(validationStatusStr)
		account.Metadata, _ = jsonbToMap(metadataJSON)

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *externalAccountRepository) GetByAccountDetails(ctx context.Context, accountNumber, routingNumber, currency string) (*entities.ExternalAccount, error) {
	query := `
		SELECT id, userId, accountNumber, routingNumber, iban, swiftCode,
			bankName, accountHolderName, currency, accountType, countryCode,
			validationStatus, validatedAt, validationExpiresAt, lastValidatedAt,
			nickname, isDefault, metadata, createdAt, updatedAt
		FROM external_accounts
		WHERE accountNumber = $1 AND routingNumber = $2 AND currency = $3
		LIMIT 1
	`

	account := &entities.ExternalAccount{}
	var validationStatusStr string
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, accountNumber, routingNumber, currency).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.RoutingNumber,
		&account.IBAN,
		&account.SwiftCode,
		&account.BankName,
		&account.AccountHolderName,
		&account.Currency,
		&account.AccountType,
		&account.CountryCode,
		&validationStatusStr,
		&account.ValidatedAt,
		&account.ValidationExpiresAt,
		&account.LastValidatedAt,
		&account.Nickname,
		&account.IsDefault,
		&metadataJSON,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("external account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get external account: %w", err)
	}

	account.ValidationStatus = entities.ValidationStatus(validationStatusStr)
	account.Metadata, _ = jsonbToMap(metadataJSON)

	return account, nil
}

func (r *externalAccountRepository) Update(ctx context.Context, account *entities.ExternalAccount) error {
	query := `
		UPDATE external_accounts
		SET validationStatus = $2, validatedAt = $3, validationExpiresAt = $4,
			lastValidatedAt = $5, nickname = $6, isDefault = $7, metadata = $8, updatedAt = $9
		WHERE id = $1
	`

	metadataJSON, _ := mapToJSONB(account.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		string(account.ValidationStatus),
		account.ValidatedAt,
		account.ValidationExpiresAt,
		account.LastValidatedAt,
		account.Nickname,
		account.IsDefault,
		metadataJSON,
		account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update external account: %w", err)
	}
	return nil
}

func (r *externalAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM external_accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete external account: %w", err)
	}
	return nil
}
