package repositories

import (
	"context"
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entities.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	GetByPaymentReference(ctx context.Context, reference string) (*entities.Payment, error)
	GetByRequestID(ctx context.Context, requestID uuid.UUID) (*entities.Payment, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Payment, error)
	Update(ctx context.Context, payment *entities.Payment) error
}
