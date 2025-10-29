package service

import (
	"context"
	"fmt"

	"github.com/chandra-shekhar/internal-transfers/internal/errs"
	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/chandra-shekhar/internal-transfers/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	accountRepo repository.AccountRepository
	logger      *zerolog.Logger
}

func NewAccountService(accountRepo repository.AccountRepository, logger *zerolog.Logger) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *model.CreateAccountRequest) (*model.Account, error) {
	// Parse the balance string to decimal
	balance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return nil, fmt.Errorf("invalid balance format: %w", err)
	}

	// Check if balance is negative
	if balance.IsNegative() {
		return nil, errs.ErrInvalidBalance
	}

	// Check if balance exceeds maximum allowed
	// NUMERIC(20,5) means max is 999999999999999.99999
	maxBalance, _ := decimal.NewFromString("999999999999999.99999")
	if balance.GreaterThan(maxBalance) {
		return nil, errs.WrapHTTPError(errs.ErrBalanceOverflow, "initial balance exceeds maximum allowed")
	}

	// Create the account
	account, err := s.accountRepo.Create(ctx, req.AccountID, balance)
	if err != nil {
		// Check if it's a duplicate key error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // PostgreSQL unique violation code
			return nil, errs.WrapHTTPError(errs.ErrAccountExists, "account with ID %d already exists", req.AccountID)
		}
		s.logger.Error().Err(err).Int64("account_id", req.AccountID).Msg("failed to create account")
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	s.logger.Info().
		Int64("account_id", account.ID).
		Str("balance", account.Balance.String()).
		Msg("account created successfully")

	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (*model.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		if err.Error() == "account not found" {
			return nil, errs.WrapHTTPError(errs.ErrAccountNotFound, "account with ID %d not found", accountID)
		}
		s.logger.Error().Err(err).Int64("account_id", accountID).Msg("failed to get account")
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &model.AccountResponse{
		AccountID: account.ID,
		Balance:   account.Balance.String(),
	}, nil
}

// ValidateAccountExists checks if an account exists
func (s *AccountService) ValidateAccountExists(ctx context.Context, accountID int64) error {
	_, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		if err.Error() == "account not found" {
			return fmt.Errorf("account with ID %d not found", accountID)
		}
		return fmt.Errorf("failed to validate account: %w", err)
	}
	return nil
}
