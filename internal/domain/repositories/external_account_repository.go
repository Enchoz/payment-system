package repositories

import (
	"context"
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
)

type ExternalAccountRepository interface {
	Create(ctx context.Context, account *entities.ExternalAccount) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ExternalAccount, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ExternalAccount, error)
	GetByAccountDetails(ctx context.Context, accountNumber, routingNumber, currency string) (*entities.ExternalAccount, error)
	Update(ctx context.Context, account *entities.ExternalAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
}
