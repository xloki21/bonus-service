package account

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/repo"
	"github.com/xloki21/bonus-service/pkg/log"
)

type Service struct {
	accounts repo.Account
}

func (a *Service) Credit(ctx context.Context, id account.UserID, value uint) error {
	return a.accounts.Credit(ctx, id, value)
}

func (a *Service) Debit(ctx context.Context, id account.UserID, value uint) error {
	return a.accounts.Debit(ctx, id, value)
}

func (a *Service) CreateAccount(ctx context.Context, account account.Account) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}

	if err := a.accounts.Create(ctx, account); err != nil {
		logger.Warnf("account creation failed: %s", err.Error())
		return err
	}
	return nil
}

func NewAccountService(accounts repo.Account) *Service {
	return &Service{accounts: accounts}
}
