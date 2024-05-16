package service

import (
	"context"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/client/accrual"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repo"
	accountSvc "github.com/xloki21/bonus-service/internal/service/account"
	orderSvc "github.com/xloki21/bonus-service/internal/service/order"
	transactionSvc "github.com/xloki21/bonus-service/internal/service/transaction"
)

type Order interface {
	Register(context.Context, order.Order) error
}

type Account interface {
	CreateAccount(context.Context, account.Account) error
	Credit(context.Context, string, uint) error
	Debit(context.Context, string, uint) error
}

type Transaction interface {
	Polling(ctx context.Context) error
}

type Service struct {
	Account
	Order
	Transaction
}

func NewService(repos *repo.Repository, cfg config.AppConfig) *Service {
	accrualClient := accrual.NewClient(cfg.AccrualService)
	return &Service{
		Account:     accountSvc.NewAccountService(repos.Account),
		Order:       orderSvc.NewOrderService(repos.Order),
		Transaction: transactionSvc.NewTransactionService(repos.Transaction, accrualClient, cfg.TransactionServiceConfig),
	}
}
