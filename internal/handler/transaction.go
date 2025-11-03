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
			return h.RespondWithHTTPError(c, httpErr)
		}

		// Handle other specific errors
		if strings.Contains(err.Error(), "invalid amount format") {
			return h.RespondWithHTTPError(c, errs.ErrInvalidFormat.WithMessage("Invalid amount format"))
		}

		h.Logger.Error().Err(err).Str("error_message", err.Error()).Msg("failed to create transaction")
		return h.RespondWithHTTPError(c, errs.ErrInternalError.WithMessage("Failed to process transaction"))
	}

	// Return empty response on success as per requirement
	return c.NoContent(http.StatusCreated)
}
