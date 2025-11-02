package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// Common business errors for the internal transfers system
var (
	ErrInsufficientBalance = &HTTPError{
		Code:     "INSUFFICIENT_BALANCE",
		Message:  "Insufficient balance in source account",
		Status:   http.StatusBadRequest,
		Override: false,
	}

	ErrAccountNotFound = &HTTPError{
		Code:     "ACCOUNT_NOT_FOUND",
		Message:  "Account not found",
		Status:   http.StatusNotFound,
		Override: false,
	}

	ErrSourceAccountNotFound = &HTTPError{
		Code:     "SOURCE_NOT_FOUND",
		Message:  "Source account not found",
		Status:   http.StatusNotFound,
		Override: false,
	}

	ErrDestinationAccountNotFound = &HTTPError{
		Code:     "DESTINATION_NOT_FOUND",
		Message:  "Destination account not found",
		Status:   http.StatusNotFound,
		Override: false,
	}

	ErrInvalidAmount = &HTTPError{
		Code:     "INVALID_AMOUNT",
		Message:  "Invalid amount",
		Status:   http.StatusBadRequest,
		Override: false,
	}

	ErrAmountMustBePositive = &HTTPError{
		Code:     "INVALID_AMOUNT",
		Message:  "Amount must be positive",
		Status:   http.StatusBadRequest,
		Override: false,
	}

	ErrSameAccount = &HTTPError{
		Code:     "SAME_ACCOUNT",
		Message:  "Source and destination accounts must be different",
		Status:   http.StatusBadRequest,
		Override: false,
	}

	ErrAccountExists = &HTTPError{
		Code:     "ACCOUNT_EXISTS",
		Message:  "Account already exists",
		Status:   http.StatusConflict,
		Override: false,
	}

	ErrInvalidBalance = &HTTPError{
		Code:     "INVALID_BALANCE",
		Message:  "Initial balance cannot be negative",
		Status:   http.StatusBadRequest,
		Override: false,
	}

	ErrBalanceOverflow = &HTTPError{
		Code:     "BALANCE_OVERFLOW",
		Message:  "Transaction would exceed maximum account balance",
		Status:   http.StatusBadRequest,
		Override: false,
	}
)

// IsHTTPError checks if an error is an HTTPError
func IsHTTPError(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr, true
	}
	return nil, false
}

// WrapHTTPError wraps an HTTPError with additional context
func WrapHTTPError(httpErr *HTTPError, format string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code:     httpErr.Code,
		Message:  fmt.Sprintf(format, args...),
		Status:   httpErr.Status,
		Override: httpErr.Override,
	}
}
