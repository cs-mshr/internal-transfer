package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/chandra-shekhar/internal-transfers/internal/errs"
	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/chandra-shekhar/internal-transfers/internal/service"
	"github.com/labstack/echo/v4"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	*BaseHandler
	accountService *service.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(base *BaseHandler, accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		BaseHandler:    base,
		accountService: accountService,
	}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	var req model.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleBindError(c, err)
	}

	if err := c.Validate(req); err != nil {
		return h.HandleValidationError(c, err)
	}

	_, err := h.accountService.CreateAccount(c.Request().Context(), &req)
	if err != nil {
		// Check for HTTPError first
		if httpErr, ok := errs.IsHTTPError(err); ok {
			return h.RespondError(c, httpErr.Status, httpErr.Code, httpErr.Message)
		}

		// Handle other specific errors
		if strings.Contains(err.Error(), "invalid balance format") {
			return h.RespondError(c, http.StatusBadRequest, "INVALID_FORMAT", "Invalid balance format")
		}

		h.Logger.Error().Err(err).Msg("failed to create account")
		return h.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create account")
	}

	// Return empty response on success as per requirement
	return c.NoContent(http.StatusCreated)
}

// GetAccount handles GET /accounts/{account_id}
func (h *AccountHandler) GetAccount(c echo.Context) error {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		return h.RespondError(c, http.StatusBadRequest, "INVALID_ACCOUNT_ID", "Invalid account ID format")
	}

	response, err := h.accountService.GetAccount(c.Request().Context(), accountID)
	if err != nil {
		// Check for HTTPError first
		if httpErr, ok := errs.IsHTTPError(err); ok {
			return h.RespondError(c, httpErr.Status, httpErr.Code, httpErr.Message)
		}

		h.Logger.Error().Err(err).Msg("failed to get account")
		return h.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get account")
	}

	return h.RespondOK(c, response)
}
