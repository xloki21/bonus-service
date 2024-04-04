package account

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/account"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=account
type accountRepository interface {
	Create(context.Context, account.Account) error
	Delete(context.Context, account.Account) error
	FindByID(context.Context, account.UserID) (*account.Account, error)
	GetBalance(context.Context, account.UserID) (int, error)
	Credit(context.Context, account.UserID, uint) error
	Debit(context.Context, account.UserID, uint) error
}
