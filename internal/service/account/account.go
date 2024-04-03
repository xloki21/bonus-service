package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/pkg/log"
)

type Account interface {
	CreateAccount(context.Context, int) (*account.Account, error)
	Credit(context.Context, account.UserID, int) error
	Debit(context.Context, account.UserID, int) error
}

type AccountService struct {
	accounts repository.Account
	logger   log.Logger
}

func (a *AccountService) Credit(ctx context.Context, id account.UserID, value int) error {

	if value < 0 {
		a.logger.Warnf("credit value is negative: %d", value)
		return apperr.InvalidCreditValue
	}

	return a.accounts.Credit(ctx, id, value)
}

func (a *AccountService) Debit(ctx context.Context, id account.UserID, value int) error {
	if value < 0 {
		a.logger.Warnf("debit value is negative: %d", value)
		return apperr.InvalidDebitValue
	}

	return a.accounts.Debit(ctx, id, value)
}

func (a *AccountService) CreateAccount(ctx context.Context, value int) (*account.Account, error) {

	acc := &account.Account{
		ID:      account.UserID(uuid.NewString()),
		Balance: value,
	}

	if err := acc.Validate(); err != nil {
		a.logger.Warnf("account validation failed: %s", err.Error())
		return nil, err
	}

	err := a.accounts.Create(ctx, acc)
	if err != nil {
		a.logger.Warnf("account creation failed: %s", err.Error())
		return nil, err
	}
	return acc, nil
}

func NewAccountService(accounts repository.Account, logger log.Logger) *AccountService {
	return &AccountService{accounts: accounts, logger: logger}
}
