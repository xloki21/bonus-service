package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/repository"
)

type AccountService struct {
	accounts repository.Account
}

func (a *AccountService) Credit(ctx context.Context, id account.UserID, value int) error {

	if value < 0 {
		return apperr.InvalidCreditValue
	}

	return a.accounts.Credit(ctx, id, value)
}

func (a *AccountService) Debit(ctx context.Context, id account.UserID, value int) error {
	if value < 0 {
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
		return nil, err
	}

	err := a.accounts.Create(ctx, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func NewAccountService(accounts repository.Account) *AccountService {
	return &AccountService{accounts: accounts}
}
