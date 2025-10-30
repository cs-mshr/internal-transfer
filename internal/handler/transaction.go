package handler

import (
	"net/http"

	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/chandra-shekhar/internal-transfers/internal/service"
	"github.com/labstack/echo/v4"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	*BaseHandler
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(base *BaseHandler, transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		BaseHandler:        base,
		transactionService: transactionService,
	}
}

// CreateTransaction handles POST /transactions
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var req model.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleBindError(c, err)
	}

	if err := c.Validate(req); err != nil {
		return h.HandleValidationError(c, err)
	}

	_, err := h.transactionService.CreateTransaction(c.Request().Context(), &req)
	if err != nil {
		// Check for specific errors
		if err.Error() == "amount must be positive" {
			return h.RespondError(c, http.StatusBadRequest, "INVALID_AMOUNT", err.Error())
		}
		if err.Error() == "source and destination accounts must be different" {
			return h.RespondError(c, http.StatusBadRequest, "SAME_ACCOUNT", err.Error())
		}
		if contains(err.Error(), "invalid amount format") {
			return h.RespondError(c, http.StatusBadRequest, "INVALID_FORMAT", "Invalid amount format")
		}
		if contains(err.Error(), "source account not found") || contains(err.Error(), "account with ID") && contains(err.Error(), "not found") && contains(err.Error(), "source") {
			return h.RespondError(c, http.StatusNotFound, "SOURCE_NOT_FOUND", "Source account not found")
		}
		if contains(err.Error(), "destination account not found") || contains(err.Error(), "account with ID") && contains(err.Error(), "not found") && contains(err.Error(), "destination") {
			return h.RespondError(c, http.StatusNotFound, "DESTINATION_NOT_FOUND", "Destination account not found")
		}
		if contains(err.Error(), "account with ID") && contains(err.Error(), "not found") {
			return h.RespondError(c, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "Account not found")
		}
		if contains(err.Error(), "insufficient balance") {
			return h.RespondError(c, http.StatusBadRequest, "INSUFFICIENT_BALANCE", "Insufficient balance in source account")
		}

		h.Logger.Error().Err(err).Msg("failed to create transaction")
		return h.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process transaction")
	}

	// Return empty response on success as per requirement
	return c.NoContent(http.StatusCreated)
}
