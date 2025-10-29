package handler

import (
	"github.com/chandra-shekhar/internal-transfers/internal/server"
	"github.com/chandra-shekhar/internal-transfers/internal/service"
)

type Handlers struct {
	Health      *HealthHandler
	OpenAPI     *OpenAPIHandler
	Account     *AccountHandler
	Transaction *TransactionHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	base := NewBaseHandler(s)

	return &Handlers{
		Health:      NewHealthHandler(s),
		OpenAPI:     NewOpenAPIHandler(s),
		Account:     NewAccountHandler(base, services.Account),
		Transaction: NewTransactionHandler(base, services.Transaction),
	}
}
