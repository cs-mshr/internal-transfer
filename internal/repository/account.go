package repository

import (
	"context"
	"fmt"

	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/chandra-shekhar/internal-transfers/internal/server"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// AccountRepository handles account-related database operations
type AccountRepository struct {
	db *pgxpool.Pool
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(s *server.Server) *AccountRepository {
	return &AccountRepository{
		db: s.DB.Pool,
	}
}

// Create creates a new account with the given ID and initial balance
func (r *AccountRepository) Create(ctx context.Context, accountID int64, initialBalance decimal.Decimal) (*model.Account, error) {
	query := `
		INSERT INTO accounts (id, balance)
		VALUES ($1, $2)
		RETURNING id, balance, created_at, updated_at
	`

	var account model.Account
	err := r.db.QueryRow(ctx, query, accountID, initialBalance).Scan(
		&account.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return &account, nil
}

// GetByID retrieves an account by its ID
func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	query := `
		SELECT id, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	var account model.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// UpdateBalance updates the balance of an account
// This should be called within a transaction for consistency
func (r *AccountRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, id int64, newBalance decimal.Decimal) error {
	query := `
		UPDATE accounts
		SET balance = $2
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, query, id, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// GetByIDForUpdate retrieves an account by its ID with a row lock for update
// This should be used within a transaction to prevent concurrent updates
func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id int64) (*model.Account, error) {
	query := `
		SELECT id, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	var account model.Account
	err := tx.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account for update: %w", err)
	}

	return &account, nil
}

// BeginTx starts a new database transaction
func (r *AccountRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
