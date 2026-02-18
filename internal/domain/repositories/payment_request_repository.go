package repositories

import (
	"context"
	"payment-system/internal/domain/entities"

	"github.com/google/uuid"
)

type PaymentRequestRepository interface {
	Create(ctx context.Context, request *entities.PaymentRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PaymentRequest, error)
	GetByRequestReference(ctx context.Context, reference string) (*entities.PaymentRequest, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*entities.PaymentRequest, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.PaymentRequest, error)
	Update(ctx context.Context, request *entities.PaymentRequest) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.PaymentRequestStatus) error
}
