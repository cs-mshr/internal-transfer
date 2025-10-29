package repository

import (
	"context"

	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// AccountRepository defines the interface for account-related database operations
type AccountRepository interface {
	Create(ctx context.Context, accountID int64, initialBalance decimal.Decimal) (*model.Account, error)
	GetByID(ctx context.Context, id int64) (*model.Account, error)
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id int64) (*model.Account, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, id int64, newBalance decimal.Decimal) error
}

// TransactionRepository defines the interface for transaction-related database operations
type TransactionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, transaction *model.Transaction) error
	UpdateStatus(ctx context.Context, tx pgx.Tx, id int64, status model.TransactionStatus) error
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
}
