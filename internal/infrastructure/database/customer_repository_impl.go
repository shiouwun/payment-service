package database

import (
	"context"
	"database/sql"

	"github.com/company/payment-service/internal/domain/entity"
	"github.com/company/payment-service/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type customerRepositoryImpl struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) repository.CustomerRepository {
	return &customerRepositoryImpl{
		db: db,
	}
}

func (r *customerRepositoryImpl) Create(ctx context.Context, customer *entity.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	return err
}

func (r *customerRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	var customer entity.Customer
	query := `
		SELECT id, name, email, phone, created_at, updated_at
		FROM customers
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &customer, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &customer, err
}

func (r *customerRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	var customer entity.Customer
	query := `
		SELECT id, name, email, phone, created_at, updated_at
		FROM customers
		WHERE email = $1
	`
	err := r.db.GetContext(ctx, &customer, query, email)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &customer, err
}

func (r *customerRepositoryImpl) Update(ctx context.Context, customer *entity.Customer) error {
	query := `
		UPDATE customers
		SET name = $2, email = $3, phone = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.UpdatedAt,
	)
	return err
}

func (r *customerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}