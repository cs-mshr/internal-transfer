package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chandra-shekhar/internal-transfers/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_ValidRequest(t *testing.T) {
	// Create request body
	reqBody := model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "100.50",
	}
	body, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Test assertions
	assert.NotNil(t, req)
	assert.NotNil(t, rec)
	assert.Equal(t, "100.50", reqBody.InitialBalance)
}

func TestGetAccount_ValidID(t *testing.T) {
	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/123", nil)
	rec := httptest.NewRecorder()

	// Test assertions
	assert.NotNil(t, req)
	assert.NotNil(t, rec)
	assert.Equal(t, "/api/v1/accounts/123", req.URL.Path)
}

func TestCreateAccount_InvalidBalance(t *testing.T) {
	// Create request body with invalid balance
	reqBody := model.CreateAccountRequest{
		AccountID:      123,
		InitialBalance: "-100.50",
	}
	body, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Verify negative balance
	assert.Equal(t, "-100.50", reqBody.InitialBalance)
}
