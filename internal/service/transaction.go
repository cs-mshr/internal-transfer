package service

import (
	"context"
	"fmt"

	"github.com/chandra-shekhar/internal-transfers/internal/database"
	"github.com/chandra-shekhar/internal-transfers/internal/errs"
	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/chandra-shekhar/internal-transfers/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type TransactionService struct {
	db              database.DB
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	logger          *zerolog.Logger
}

func NewTransactionService(db database.DB, accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		logger:          logger,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	// Parse the amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate amount
	if amount.IsNegative() || amount.IsZero() {
		return nil, errs.ErrAmountMustBePositive
	}

	// Validate that source and destination are different
	if req.SourceAccountID == req.DestinationAccountID {
		return nil, errs.ErrSameAccount
	}

	// Verify accounts exist before starting transaction
	_, err = s.accountRepo.GetByID(ctx, req.SourceAccountID)
	if err != nil {
		if err.Error() == "account not found" {
			return nil, errs.ErrSourceAccountNotFound
		}
		return nil, fmt.Errorf("failed to verify source account: %w", err)
	}

	_, err = s.accountRepo.GetByID(ctx, req.DestinationAccountID)
	if err != nil {
		if err.Error() == "account not found" {
			return nil, errs.ErrDestinationAccountNotFound
		}
		return nil, fmt.Errorf("failed to verify destination account: %w", err)
	}

	// Start a database transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				s.logger.Error().Err(rollbackErr).Msg("failed to rollback transaction")
			}
		}
	}()

	// Create transaction record
	transaction := &model.Transaction{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               amount,
		Status:               model.TransactionStatusPending,
	}
	err = s.transactionRepo.Create(ctx, tx, transaction)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create transaction record")
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Lock and get source account (with FOR UPDATE to prevent concurrent modifications)
	sourceAccount, err := s.accountRepo.GetByIDForUpdate(ctx, tx, req.SourceAccountID)
	if err != nil {
		if err.Error() == "account not found" {
			updateErr := s.transactionRepo.UpdateStatus(ctx, tx, transaction.ID, model.TransactionStatusFailed)
			if updateErr != nil {
				s.logger.Error().Err(updateErr).Msg("failed to update transaction status")
			}
			return nil, errs.ErrSourceAccountNotFound
		}
		return nil, fmt.Errorf("failed to get source account: %w", err)
	}

	// Lock and get destination account
	destAccount, err := s.accountRepo.GetByIDForUpdate(ctx, tx, req.DestinationAccountID)
	if err != nil {
		if err.Error() == "account not found" {
			updateErr := s.transactionRepo.UpdateStatus(ctx, tx, transaction.ID, model.TransactionStatusFailed)
			if updateErr != nil {
				s.logger.Error().Err(updateErr).Msg("failed to update transaction status")
			}
			return nil, errs.ErrDestinationAccountNotFound
		}
		return nil, fmt.Errorf("failed to get destination account: %w", err)
	}

	// Check if source account has sufficient balance
	if sourceAccount.Balance.LessThan(amount) {
		updateErr := s.transactionRepo.UpdateStatus(ctx, tx, transaction.ID, model.TransactionStatusFailed)
		if updateErr != nil {
			s.logger.Error().Err(updateErr).Msg("failed to update transaction status")
		}
		return nil, errs.ErrInsufficientBalance
	}

	// Calculate new balances
	newSourceBalance := sourceAccount.Balance.Sub(amount)
	newDestBalance := destAccount.Balance.Add(amount)

	// Update source account balance
	if err := s.accountRepo.UpdateBalance(ctx, tx, req.SourceAccountID, newSourceBalance); err != nil {
		s.logger.Error().Err(err).Msg("failed to update source account balance")
		return nil, fmt.Errorf("failed to update source account balance: %w", err)
	}

	// Update destination account balance
	if err := s.accountRepo.UpdateBalance(ctx, tx, req.DestinationAccountID, newDestBalance); err != nil {
		s.logger.Error().Err(err).Msg("failed to update destination account balance")
		return nil, fmt.Errorf("failed to update destination account balance: %w", err)
	}

	// Update transaction status to completed
	if err := s.transactionRepo.UpdateStatus(ctx, tx, transaction.ID, model.TransactionStatusCompleted); err != nil {
		s.logger.Error().Err(err).Msg("failed to update transaction status")
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Int64("transaction_id", transaction.ID).
		Int64("source_account_id", req.SourceAccountID).
		Int64("destination_account_id", req.DestinationAccountID).
		Str("amount", amount.String()).
		Msg("transaction completed successfully")

	return &model.TransactionResponse{
		ID:                   transaction.ID,
		SourceAccountID:      transaction.SourceAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Amount:               transaction.Amount.String(),
		Status:               model.TransactionStatusCompleted,
		CreatedAt:            transaction.CreatedAt,
	}, nil
}

// GetTransaction retrieves a transaction by its ID
func (s *TransactionService) GetTransaction(ctx context.Context, transactionID int64) (*model.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		if err.Error() == "transaction not found" {
			return nil, fmt.Errorf("transaction with ID %d not found", transactionID)
		}
		s.logger.Error().Err(err).Int64("transaction_id", transactionID).Msg("failed to get transaction")
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &model.TransactionResponse{
		ID:                   transaction.ID,
		SourceAccountID:      transaction.SourceAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Amount:               transaction.Amount.String(),
		Status:               transaction.Status,
		CreatedAt:            transaction.CreatedAt,
	}, nil
}

// processTransfer handles the actual money transfer logic within a transaction
func (s *TransactionService) processTransfer(ctx context.Context, tx pgx.Tx, sourceID, destID int64, amount decimal.Decimal) error {
	// This is a helper method that could be expanded with additional business logic
	// For now, the logic is in CreateTransaction method
	return nil
}
