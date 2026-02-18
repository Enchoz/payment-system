package repositories

import (
	"context"
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
)

type LedgerRepository interface {
	Create(ctx context.Context, entry *entities.LedgerEntry) error
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*entities.LedgerEntry, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entities.LedgerEntry, error)
}
