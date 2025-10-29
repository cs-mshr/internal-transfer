package service

import (
	"github.com/chandra-shekhar/internal-transfers/internal/repository"
	"github.com/chandra-shekhar/internal-transfers/internal/server"
)

type Services struct {
	Account     *AccountService
	Transaction *TransactionService
}

func NewServices(s *server.Server, repos *repository.Repositories) *Services {
	return &Services{
		Account:     NewAccountService(repos.Account, s.Logger),
		Transaction: NewTransactionService(repos.Account, repos.Transaction, s.Logger),
	}
}
