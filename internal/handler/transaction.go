package handler

import (
	"net/http"
	"strings"

	"github.com/chandra-shekhar/internal-transfers/internal/errs"
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
		// Check for HTTPError first
		if httpErr, ok := errs.IsHTTPError(err); ok {
			return h.RespondError(c, httpErr.Status, httpErr.Code, httpErr.Message)
		}

		// Handle other specific errors
		if strings.Contains(err.Error(), "invalid amount format") {
			return h.RespondError(c, http.StatusBadRequest, "INVALID_FORMAT", "Invalid amount format")
		}

		h.Logger.Error().Err(err).Str("error_message", err.Error()).Msg("failed to create transaction")
		return h.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process transaction")
	}

	// Return empty response on success as per requirement
	return c.NoContent(http.StatusCreated)
}
