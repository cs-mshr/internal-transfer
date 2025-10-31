package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chandra-shekhar/internal-transfers/internal/database"
	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/chandra-shekhar/internal-transfers/internal/server"
	"github.com/jackc/pgx/v5"
)

type transactionRepository struct {
	db database.DB
}

func NewTransactionRepository(s *server.Server) TransactionRepository {
	return &transactionRepository{
		db: s.DB,
	}
}

func (r *transactionRepository) Create(ctx context.Context, tx pgx.Tx, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions (source_account_id, destination_account_id, amount, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, source_account_id, destination_account_id, amount, status, created_at, completed_at
	`

	err := tx.QueryRow(ctx, query, transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount, model.TransactionStatusPending).Scan(
		&transaction.ID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id int64, status model.TransactionStatus) error {
	var query string
	var args []interface{}

	if status == model.TransactionStatusCompleted {
		query = `
			UPDATE transactions
			SET status = $2, completed_at = $3
			WHERE id = $1
		`
		args = []interface{}{id, status, time.Now()}
	} else {
		query = `
			UPDATE transactions
			SET status = $2
			WHERE id = $1
		`
		args = []interface{}{id, status}
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	query := `
		SELECT id, source_account_id, destination_account_id, amount, status, created_at, completed_at
		FROM transactions
		WHERE id = $1
	`

	var transaction model.Transaction
	err := r.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.SourceAccountID,
		&transaction.DestinationAccountID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// GetByAccountID retrieves all transactions for a specific account
func (r *transactionRepository) GetByAccountID(ctx context.Context, accountID int64, limit, offset int) ([]*model.Transaction, error) {
	query := `
		SELECT id, source_account_id, destination_account_id, amount, status, created_at, completed_at
		FROM transactions
		WHERE source_account_id = $1 OR destination_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SourceAccountID,
			&transaction.DestinationAccountID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}
