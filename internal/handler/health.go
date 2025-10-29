package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chandra-shekhar/internal-transfers/internal/middleware"
	"github.com/chandra-shekhar/internal-transfers/internal/server"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	server *server.Server
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		server: s,
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	start := time.Now()
	logger := middleware.GetLogger(c).With().
		Str("operation", "health_check").
		Logger()

	logger.Debug().Msg("checking health status")

	// Check database connectivity
	dbStatus := h.checkDatabaseHealth()

	overallStatus := "healthy"
	if dbStatus.Status != "healthy" {
		overallStatus = "degraded"
	}

	health := HealthResponse{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		Version:     h.server.Config.Primary.Env,
		ServiceName: "internal-transfers",
		Checks: map[string]ComponentHealth{
			"database": dbStatus,
		},
	}

	totalDuration := time.Since(start)
	logger.Info().
		Str("status", overallStatus).
		Dur("duration", totalDuration).
		Msg("health check completed")

	// Return appropriate status code based on health
	statusCode := http.StatusOK
	if overallStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, health)
}

func (h *HealthHandler) checkDatabaseHealth() ComponentHealth {
	start := time.Now()
	status := "healthy"
	message := "connected"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.DB.Pool.Ping(ctx); err != nil {
		status = "unhealthy"
		message = fmt.Sprintf("connection failed: %s", err.Error())
	}

	return ComponentHealth{
		Status:   status,
		Message:  message,
		Duration: time.Since(start).Milliseconds(),
	}
}

// Health response types

type HealthResponse struct {
	Status      string                     `json:"status"`
	Timestamp   time.Time                  `json:"timestamp"`
	Version     string                     `json:"version"`
	ServiceName string                     `json:"service_name"`
	Checks      map[string]ComponentHealth `json:"checks"`
}

type ComponentHealth struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	Duration int64  `json:"duration_ms,omitempty"`
}
