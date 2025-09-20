package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/company/payment-service/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByReference(ctx context.Context, reference string) (*entity.Payment, error) {
	args := m.Called(ctx, reference)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
	args := m.Called(ctx, merchantID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
	args := m.Called(ctx, customerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Payment), args.Error(1)
}

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) Create(ctx context.Context, merchant *entity.Merchant) error {
	args := m.Called(ctx, merchant)
	return args.Error(0)
}

func (m *MockMerchantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) GetByAPIKey(ctx context.Context, apiKey string) (*entity.Merchant, error) {
	args := m.Called(ctx, apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) Update(ctx context.Context, merchant *entity.Merchant) error {
	args := m.Called(ctx, merchant)
	return args.Error(0)
}

func (m *MockMerchantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPaymentUseCase_CreatePayment(t *testing.T) {
	ctx := context.Background()

	merchantID := uuid.New()
	customerID := uuid.New()

	tests := []struct {
		name          string
		request       CreatePaymentRequest
		setupMocks    func(*MockPaymentRepository, *MockMerchantRepository, *MockCustomerRepository)
		expectedError string
	}{
		{
			name: "successful payment creation",
			request: CreatePaymentRequest{
				MerchantID:  merchantID,
				CustomerID:  customerID,
				Amount:      10000, // $100.00
				Currency:    "USD",
				Method:      entity.PaymentMethodCreditCard,
				Description: "Test payment",
				Reference:   "REF001",
			},
			setupMocks: func(paymentRepo *MockPaymentRepository, merchantRepo *MockMerchantRepository, customerRepo *MockCustomerRepository) {
				merchant := &entity.Merchant{
					ID:       merchantID,
					Name:     "Test Merchant",
					IsActive: true,
				}
				customer := &entity.Customer{
					ID:   customerID,
					Name: "Test Customer",
				}

				merchantRepo.On("GetByID", ctx, merchantID).Return(merchant, nil)
				customerRepo.On("GetByID", ctx, customerID).Return(customer, nil)
				paymentRepo.On("Create", ctx, mock.AnythingOfType("*entity.Payment")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "inactive merchant",
			request: CreatePaymentRequest{
				MerchantID:  merchantID,
				CustomerID:  customerID,
				Amount:      10000,
				Currency:    "USD",
				Method:      entity.PaymentMethodCreditCard,
				Description: "Test payment",
			},
			setupMocks: func(paymentRepo *MockPaymentRepository, merchantRepo *MockMerchantRepository, customerRepo *MockCustomerRepository) {
				merchant := &entity.Merchant{
					ID:       merchantID,
					Name:     "Test Merchant",
					IsActive: false,
				}

				merchantRepo.On("GetByID", ctx, merchantID).Return(merchant, nil)
			},
			expectedError: "merchant is not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepo := new(MockPaymentRepository)
			merchantRepo := new(MockMerchantRepository)
			customerRepo := new(MockCustomerRepository)

			tt.setupMocks(paymentRepo, merchantRepo, customerRepo)

			useCase := NewPaymentUseCase(paymentRepo, merchantRepo, customerRepo)

			payment, err := useCase.CreatePayment(ctx, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)
				assert.Equal(t, tt.request.MerchantID, payment.MerchantID)
				assert.Equal(t, tt.request.CustomerID, payment.CustomerID)
				assert.Equal(t, tt.request.Amount, payment.Amount)
				assert.Equal(t, entity.PaymentStatusPending, payment.Status)
			}

			paymentRepo.AssertExpectations(t)
			merchantRepo.AssertExpectations(t)
			customerRepo.AssertExpectations(t)
		})
	}
}

func TestPaymentUseCase_ProcessPayment(t *testing.T) {
	ctx := context.Background()
	paymentID := uuid.New()

	tests := []struct {
		name          string
		paymentID     uuid.UUID
		setupMocks    func(*MockPaymentRepository)
		expectedError string
	}{
		{
			name:      "successful payment processing",
			paymentID: paymentID,
			setupMocks: func(paymentRepo *MockPaymentRepository) {
				payment := &entity.Payment{
					ID:     paymentID,
					Status: entity.PaymentStatusPending,
				}
				paymentRepo.On("GetByID", ctx, paymentID).Return(payment, nil)
				paymentRepo.On("UpdateStatus", ctx, paymentID, entity.PaymentStatusCompleted).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "payment already completed",
			paymentID: paymentID,
			setupMocks: func(paymentRepo *MockPaymentRepository) {
				payment := &entity.Payment{
					ID:     paymentID,
					Status: entity.PaymentStatusCompleted,
				}
				paymentRepo.On("GetByID", ctx, paymentID).Return(payment, nil)
			},
			expectedError: "payment status is completed, cannot process",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepo := new(MockPaymentRepository)
			merchantRepo := new(MockMerchantRepository)
			customerRepo := new(MockCustomerRepository)

			tt.setupMocks(paymentRepo)

			useCase := NewPaymentUseCase(paymentRepo, merchantRepo, customerRepo)

			err := useCase.ProcessPayment(ctx, tt.paymentID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			paymentRepo.AssertExpectations(t)
		})
	}
}