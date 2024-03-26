package service

import (
	"context"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repository"
)

type Account interface {
	CreateAccount(context.Context, int) (*account.Account, error)
	Credit(context.Context, account.UserID, int) error
	Debit(context.Context, account.UserID, int) error
}

type Order interface {
	Register(context.Context, *order.Order) error
}

type Accrual interface {
	RequestOrderReward(ctx context.Context, order *order.Order) (uint, error)
}

type Transaction interface {
	Polling(ctx context.Context) error
}

type Service struct {
	Account
	Accrual
	Order
	Transaction
}

func NewService(repos *repository.Repository, cfg *config.AppConfig) *Service {
	return &Service{
		Account:     NewAccountService(repos.Account),
		Accrual:     NewAccrualService(repos.Transaction),
		Order:       NewOrderService(repos.Order),
		Transaction: NewTransactionService(repos.Transaction, cfg),
	}
}
