package repository

import (
	"context"

	"github.com/company/payment-service/internal/domain/entity"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	GetByReference(ctx context.Context, reference string) (*entity.Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]*entity.Payment, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*entity.Payment, error)
}

type MerchantRepository interface {
	Create(ctx context.Context, merchant *entity.Merchant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Merchant, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*entity.Merchant, error)
	Update(ctx context.Context, merchant *entity.Merchant) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
	GetByEmail(ctx context.Context, email string) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}