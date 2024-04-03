package service

import (
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/internal/service/account"
	"github.com/xloki21/bonus-service/internal/service/order"
	"github.com/xloki21/bonus-service/internal/service/transaction"
)

type Service struct {
	account.Account
	order.Order
	transaction.Transaction
}

func NewService(repos *repository.Repository, cfg config.AppConfig) *Service {
	return &Service{
		Account:     account.NewAccountService(repos.Account),
		Order:       order.NewOrderService(repos.Order),
		Transaction: transaction.NewTransactionService(repos.Transaction, cfg),
	}
}
