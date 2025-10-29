package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Account represents a bank account
type Account struct {
	ID        int64           `json:"account_id" db:"id"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	AccountID      int64  `json:"account_id" validate:"required,min=1"`
	InitialBalance string `json:"initial_balance" validate:"required,numeric"`
}

// AccountResponse represents the response for account queries
type AccountResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}
