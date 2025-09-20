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

type merchantRepository struct {
	db *sqlx.DB
}

func NewMerchantRepository(db *sqlx.DB) repository.MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) Create(ctx context.Context, merchant *entity.Merchant) error {
	query := `
		INSERT INTO merchants (id, name, email, api_key, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		merchant.ID, merchant.Name, merchant.Email, merchant.APIKey,
		merchant.IsActive, merchant.CreatedAt, merchant.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create merchant")
	}
	return nil
}

func (r *merchantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Merchant, error) {
	query := `
		SELECT id, name, email, api_key, is_active, created_at, updated_at
		FROM merchants WHERE id = $1
	`
	var merchant entity.Merchant
	err := r.db.GetContext(ctx, &merchant, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("merchant not found")
		}
		return nil, errors.Wrap(err, "failed to get merchant by id")
	}
	return &merchant, nil
}

func (r *merchantRepository) GetByAPIKey(ctx context.Context, apiKey string) (*entity.Merchant, error) {
	query := `
		SELECT id, name, email, api_key, is_active, created_at, updated_at
		FROM merchants WHERE api_key = $1
	`
	var merchant entity.Merchant
	err := r.db.GetContext(ctx, &merchant, query, apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("merchant not found")
		}
		return nil, errors.Wrap(err, "failed to get merchant by api key")
	}
	return &merchant, nil
}

func (r *merchantRepository) Update(ctx context.Context, merchant *entity.Merchant) error {
	merchant.UpdatedAt = time.Now()
	query := `
		UPDATE merchants
		SET name = $1, email = $2, api_key = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.ExecContext(ctx, query,
		merchant.Name, merchant.Email, merchant.APIKey, merchant.IsActive,
		merchant.UpdatedAt, merchant.ID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update merchant")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.New("merchant not found")
	}

	return nil
}

func (r *merchantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE merchants SET is_active = false, updated_at = $1 WHERE id = $2"
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return errors.Wrap(err, "failed to delete merchant")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get affected rows")
	}
	if rowsAffected == 0 {
		return errors.New("merchant not found")
	}

	return nil
}