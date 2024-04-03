package account

import (
	"context"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/pkg/log"
)

type Account interface {
	CreateAccount(context.Context, account.Account) error
	Credit(context.Context, account.UserID, int) error
	Debit(context.Context, account.UserID, int) error
}

type Service struct {
	accounts repository.Account
}

func (a *Service) Credit(ctx context.Context, id account.UserID, value int) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}
	if value < 0 {
		logger.Warnf("credit value is negative: %d", value)
		return apperr.InvalidCreditValue
	}

	return a.accounts.Credit(ctx, id, value)
}

func (a *Service) Debit(ctx context.Context, id account.UserID, value int) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}
	if value < 0 {
		logger.Warnf("debit value is negative: %d", value)
		return apperr.InvalidDebitValue
	}

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

func NewAccountService(accounts repository.Account) *Service {
	return &Service{accounts: accounts}
}
