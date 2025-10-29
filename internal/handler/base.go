package handler

import (
	"net/http"

	"github.com/chandra-shekhar/internal-transfers/internal/server"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// BaseHandler provides base functionality for all handlers
type BaseHandler struct {
	Server *server.Server
	Logger *zerolog.Logger
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(s *server.Server) *BaseHandler {
	return &BaseHandler{
		Server: s,
		Logger: s.Logger,
	}
}

// RespondOK sends a successful response with data
func (h *BaseHandler) RespondOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

// RespondError sends an error response
func (h *BaseHandler) RespondError(c echo.Context, statusCode int, errorCode, message string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"error": map[string]string{
			"code":    errorCode,
			"message": message,
		},
	})
}

// HandleBindError handles request binding errors
func (h *BaseHandler) HandleBindError(c echo.Context, err error) error {
	h.Logger.Error().Err(err).Msg("failed to bind request")
	return h.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
}

// HandleValidationError handles validation errors
func (h *BaseHandler) HandleValidationError(c echo.Context, err error) error {
	h.Logger.Error().Err(err).Msg("validation failed")

	// Extract validation errors
	if ve, ok := err.(validator.ValidationErrors); ok {
		// Get the first validation error
		if len(ve) > 0 {
			field := ve[0].Field()
			tag := ve[0].Tag()

			var message string
			switch tag {
			case "required":
				message = field + " is required"
			case "min":
				message = field + " is too small"
			case "max":
				message = field + " is too large"
			case "numeric":
				message = field + " must be numeric"
			case "nefield":
				message = field + " must be different from " + ve[0].Param()
			default:
				message = field + " is invalid"
			}

			return h.RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", message)
		}
	}

	return h.RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Request validation failed")
}
