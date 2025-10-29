package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction represents a money transfer between accounts
type Transaction struct {
	ID                   int64             `json:"id" db:"id"`
	SourceAccountID      int64             `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int64             `json:"destination_account_id" db:"destination_account_id"`
	Amount               decimal.Decimal   `json:"amount" db:"amount"`
	Status               TransactionStatus `json:"status" db:"status"`
	CreatedAt            time.Time         `json:"created_at" db:"created_at"`
	CompletedAt          *time.Time        `json:"completed_at,omitempty" db:"completed_at"`
}

// CreateTransactionRequest represents the request to create a new transaction
type CreateTransactionRequest struct {
	SourceAccountID      int64  `json:"source_account_id" validate:"required,min=1"`
	DestinationAccountID int64  `json:"destination_account_id" validate:"required,min=1,nefield=SourceAccountID"`
	Amount               string `json:"amount" validate:"required,numeric"`
}

// TransactionResponse represents the response for transaction creation
type TransactionResponse struct {
	ID                   int64             `json:"id"`
	SourceAccountID      int64             `json:"source_account_id"`
	DestinationAccountID int64             `json:"destination_account_id"`
	Amount               string            `json:"amount"`
	Status               TransactionStatus `json:"status"`
	CreatedAt            time.Time         `json:"created_at"`
}
