package service

import (
	"context"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/integration"
	"github.com/xloki21/bonus-service/internal/pkg/log"
	"github.com/xloki21/bonus-service/internal/repository"
	"sync"
	"time"
)

type TransactionService struct {
	cfg    *config.AppConfig
	repo   repository.Transaction
	logger log.Logger
}

// Polling is a blocking method that polls unprocessed transactions
func (t *TransactionService) Polling(ctx context.Context) error {
	t.logger.Info("polling transactions...")
	accrualServiceClient := integration.New(t.cfg.AccrualService)
	ticker := time.NewTicker(t.cfg.TransactionServiceConfig.PollingInterval)
	defer ticker.Stop() // Stop the ticker so it can be garbage collected
	for {
		select {
		case <-ctx.Done():
			t.logger.Info("polling events listener stopped")
			return ctx.Err()
		case <-ticker.C:
			t.logger.Info("polling event triggered")
			t.logger.Info("find unprocessed transactions...")
			txs, err := t.repo.FindUnprocessed(ctx, int64(t.cfg.TransactionServiceConfig.MaxTransactionsPerRequest))
			if err != nil {
				t.logger.Warnf("polling event error on find unprocessed transactions: %v", err)
				continue
			}
			t.logger.Info("processing transactions...")
			wg := &sync.WaitGroup{}
			for _, tx := range txs {
				wg.Add(1)
				go func(wg *sync.WaitGroup, tx *transaction.Transaction) {
					defer wg.Done()
					reward, err := accrualServiceClient.GetAccrual(ctx, tx)
					if err != nil {
						t.logger.Warnf("error during request to accrual service: %v", err)
						return
					}
					tx.Reward = reward
					tx.Status = transaction.PROCESSED

					if err := t.repo.Update(ctx, tx); err != nil {
						t.logger.Warnf("polling event error on update transaction: %v", err)
						return
					}
				}(wg, &tx)
			}
			wg.Wait()

			t.logger.Info("rewarding accounts...")
			if err := t.repo.RewardAccounts(ctx, int64(t.cfg.TransactionServiceConfig.MaxTransactionsPerRequest)); err != nil {
				t.logger.Warnf("polling event error on reward accounts: %v", err)
			}

		}
	}
}

func NewTransactionService(transactions repository.Transaction, cfg *config.AppConfig, logger log.Logger) *TransactionService {
	return &TransactionService{repo: transactions, cfg: cfg, logger: logger}
}
