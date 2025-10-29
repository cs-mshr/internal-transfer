package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chandra-shekhar/internal-transfers/internal/server"

	"github.com/labstack/echo/v4"
)

type OpenAPIHandler struct {
	server *server.Server
}

func NewOpenAPIHandler(s *server.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		server: s,
	}
}

func (h *OpenAPIHandler) ServeOpenAPIUI(c echo.Context) error {
	templateBytes, err := os.ReadFile("static/openapi.html")
	c.Response().Header().Set("Cache-Control", "no-cache")
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI UI template: %w", err)
	}

	templateString := string(templateBytes)
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	return c.String(http.StatusOK, templateString)
}
