package repository

import (
	"github.com/chandra-shekhar/internal-transfers/internal/server"
)

type Repositories struct {
	Account     *AccountRepository
	Transaction *TransactionRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Account:     NewAccountRepository(s),
		Transaction: NewTransactionRepository(s),
	}
}
