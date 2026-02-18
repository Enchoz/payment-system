package repositories

import (
	"context"
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
)

type PaymentHoldRepository interface {
	Create(ctx context.Context, hold *entities.PaymentHold) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PaymentHold, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entities.PaymentHold, error)
	GetByRequestID(ctx context.Context, requestID uuid.UUID) (*entities.PaymentHold, error)
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) (*entities.PaymentHold, error)
	Update(ctx context.Context, hold *entities.PaymentHold) error
	Release(ctx context.Context, holdID uuid.UUID) error
}
