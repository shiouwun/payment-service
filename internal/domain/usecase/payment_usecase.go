package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/company/payment-service/internal/domain/entity"
	"github.com/company/payment-service/internal/domain/repository"
	"github.com/company/payment-service/pkg/errors"
	"github.com/google/uuid"
)

type PaymentUseCase interface {
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*entity.Payment, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	ProcessPayment(ctx context.Context, id uuid.UUID) error
	CancelPayment(ctx context.Context, id uuid.UUID) error
	GetMerchantPayments(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]*entity.Payment, error)
}

type CreatePaymentRequest struct {
	MerchantID  uuid.UUID            `json:"merchant_id" validate:"required"`
	CustomerID  uuid.UUID            `json:"customer_id" validate:"required"`
	Amount      int64                `json:"amount" validate:"required,gt=0"`
	Currency    string               `json:"currency" validate:"required,len=3"`
	Method      entity.PaymentMethod `json:"method" validate:"required"`
	Description string               `json:"description"`
	Reference   string               `json:"reference"`
}

type paymentUseCase struct {
	paymentRepo  repository.PaymentRepository
	merchantRepo repository.MerchantRepository
	customerRepo repository.CustomerRepository
}

func NewPaymentUseCase(
	paymentRepo repository.PaymentRepository,
	merchantRepo repository.MerchantRepository,
	customerRepo repository.CustomerRepository,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo:  paymentRepo,
		merchantRepo: merchantRepo,
		customerRepo: customerRepo,
	}
}

func (uc *paymentUseCase) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*entity.Payment, error) {
	// 驗證商戶存在且活躍
	merchant, err := uc.merchantRepo.GetByID(ctx, req.MerchantID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get merchant")
	}
	if !merchant.IsActive {
		return nil, errors.New("merchant is not active")
	}

	// 驗證客戶存在
	_, err = uc.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get customer")
	}

	// 創建支付記錄
	payment := &entity.Payment{
		ID:          uuid.New(),
		MerchantID:  req.MerchantID,
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Method:      req.Method,
		Status:      entity.PaymentStatusPending,
		Description: req.Description,
		Reference:   req.Reference,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, errors.Wrap(err, "failed to create payment")
	}

	return payment, nil
}

func (uc *paymentUseCase) GetPayment(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payment")
	}
	return payment, nil
}

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, id uuid.UUID) error {
	payment, err := uc.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get payment")
	}

	if payment.Status != entity.PaymentStatusPending {
		return errors.New(fmt.Sprintf("payment status is %s, cannot process", payment.Status))
	}

	// 在實際應用中，這裡會調用第三方支付網關
	// 為了示例，我們假設支付總是成功

	if err := uc.paymentRepo.UpdateStatus(ctx, id, entity.PaymentStatusCompleted); err != nil {
		return errors.Wrap(err, "failed to update payment status")
	}

	return nil
}

func (uc *paymentUseCase) CancelPayment(ctx context.Context, id uuid.UUID) error {
	payment, err := uc.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get payment")
	}

	if payment.Status != entity.PaymentStatusPending {
		return errors.New(fmt.Sprintf("payment status is %s, cannot cancel", payment.Status))
	}

	if err := uc.paymentRepo.UpdateStatus(ctx, id, entity.PaymentStatusCancelled); err != nil {
		return errors.Wrap(err, "failed to update payment status")
	}

	return nil
}

func (uc *paymentUseCase) GetMerchantPayments(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
	payments, err := uc.paymentRepo.GetByMerchantID(ctx, merchantID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get merchant payments")
	}
	return payments, nil
}