package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodDigitalWallet PaymentMethod = "digital_wallet"
)

type Payment struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	MerchantID    uuid.UUID     `json:"merchant_id" db:"merchant_id"`
	CustomerID    uuid.UUID     `json:"customer_id" db:"customer_id"`
	Amount        int64         `json:"amount" db:"amount"` // 以分為單位避免浮點數精度問題
	Currency      string        `json:"currency" db:"currency"`
	Method        PaymentMethod `json:"method" db:"method"`
	Status        PaymentStatus `json:"status" db:"status"`
	Description   string        `json:"description" db:"description"`
	Reference     string        `json:"reference" db:"reference"` // 外部參考號
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty" db:"completed_at"`
}

type Merchant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	APIKey    string    `json:"-" db:"api_key"` // 不在JSON中暴露
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Customer struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}