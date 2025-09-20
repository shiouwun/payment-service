package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/company/payment-service/internal/domain/entity"
	"github.com/company/payment-service/internal/domain/repository"
	"github.com/company/payment-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	query := `
		INSERT INTO payments (id, merchant_id, customer_id, amount, currency, method, status, description, reference, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.MerchantID, payment.CustomerID, payment.Amount,
		payment.Currency, payment.Method, payment.Status, payment.Description,
		payment.Reference, payment.CreatedAt, payment.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create payment")
	}
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	query := `
		SELECT id, merchant_id, customer_id, amount, currency, method, status,
		       description, reference, created_at, updated_at, completed_at
		FROM payments WHERE id = $1
	`
	var payment entity.Payment
	err := r.db.GetContext(ctx, &payment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, errors.Wrap(err, "failed to get payment by id")
	}
	return &payment, nil
}

func (r *paymentRepository) GetByReference(ctx context.Context, reference string) (*entity.Payment, error) {
	query := `
		SELECT id, merchant_id, customer_id, amount, currency, method, status,
		       description, reference, created_at, updated_at, completed_at
		FROM payments WHERE reference = $1
	`
	var payment entity.Payment
	err := r.db.GetContext(ctx, &payment, query, reference)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, errors.Wrap(err, "failed to get payment by reference")
	}
	return &payment, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error {
	var completedAt *time.Time
	if status == entity.PaymentStatusCompleted {
		now := time.Now()
		completedAt = &now
	}

	query := `
		UPDATE payments
		SET status = $1, updated_at = $2, completed_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query, status, time.Now(), completedAt, id)
	if err != nil {
		return errors.Wrap(err, "failed to update payment status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.New("payment not found")
	}

	return nil
}

func (r *paymentRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
	query := `
		SELECT id, merchant_id, customer_id, amount, currency, method, status,
		       description, reference, created_at, updated_at, completed_at
		FROM payments
		WHERE merchant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var payments []*entity.Payment
	err := r.db.SelectContext(ctx, &payments, query, merchantID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payments by merchant id")
	}
	return payments, nil
}

func (r *paymentRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
	query := `
		SELECT id, merchant_id, customer_id, amount, currency, method, status,
		       description, reference, created_at, updated_at, completed_at
		FROM payments
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var payments []*entity.Payment
	err := r.db.SelectContext(ctx, &payments, query, customerID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payments by customer id")
	}
	return payments, nil
}