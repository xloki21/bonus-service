package service

import (
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/internal/service/account"
	"github.com/xloki21/bonus-service/internal/service/accrual"
	"github.com/xloki21/bonus-service/internal/service/order"
	"github.com/xloki21/bonus-service/internal/service/transaction"
	"github.com/xloki21/bonus-service/pkg/log"
)

type Service struct {
	account.Account
	accrual.Accrual
	order.Order
	transaction.Transaction
}

func NewService(repos *repository.Repository, cfg *config.AppConfig, logger log.Logger) *Service {
	return &Service{
		Account:     account.NewAccountService(repos.Account, logger),
		Accrual:     accrual.NewAccrualService(repos.Transaction, logger),
		Order:       order.NewOrderService(repos.Order, logger),
		Transaction: transaction.NewTransactionService(repos.Transaction, cfg, logger),
	}
}
