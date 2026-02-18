package repositories

import (
	"context"
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entities.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Account, error)
	GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*entities.Account, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Account, error)
	Update(ctx context.Context, account *entities.Account) error
	UpdateBalance(ctx context.Context, accountID uuid.UUID, balance, availableBalance interface{}) error
}
