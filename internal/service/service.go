package service

import (
	"context"
	"github.com/xloki21/bonus-service/config"
	a "github.com/xloki21/bonus-service/internal/entity/account"
	o "github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/internal/service/account"
	"github.com/xloki21/bonus-service/internal/service/order"
	"github.com/xloki21/bonus-service/internal/service/transaction"
)

type Order interface {
	Register(context.Context, o.Order) error
}

type Account interface {
	CreateAccount(context.Context, a.Account) error
	Credit(context.Context, a.UserID, uint) error
	Debit(context.Context, a.UserID, uint) error
}

type Transaction interface {
	Polling(ctx context.Context) error
}

type Service struct {
	Account
	Order
	Transaction
}

func NewService(repos *repository.Repository, cfg config.AppConfig) *Service {
	return &Service{
		Account:     account.NewAccountService(repos.Account),
		Order:       order.NewOrderService(repos.Order),
		Transaction: transaction.NewTransactionService(repos.Transaction, cfg),
	}
}
