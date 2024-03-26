package service

import (
	"context"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/integration"
	"github.com/xloki21/bonus-service/internal/repository"
	"log"
	"sync"
	"time"
)

const (
	//pollingInterval                   = 1 * time.Second
	maxUnprocessedTransactionPerQuery = 10
)

type TransactionService struct {
	cfg  *config.AppConfig
	repo repository.Transaction
}

// Polling is a blocking method that polls unprocessed transactions
func (t *TransactionService) Polling(ctx context.Context) error {
	accrualServiceClient := integration.New(&t.cfg.AccrualService)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(t.cfg.TransactionServiceConfig.PollingInterval):
			txs, err := t.repo.FindUnprocessed(ctx, int64(t.cfg.TransactionServiceConfig.MaxTransactionsPerRequest))
			if err != nil {
				log.Printf("polling error: %v\n", err)
				continue
			}
			wg := &sync.WaitGroup{}
			for _, tx := range txs {
				wg.Add(1)
				go func(wg *sync.WaitGroup, tx *transaction.Transaction) {
					defer wg.Done()
					reward, err := accrualServiceClient.GetAccrual(ctx, tx)
					if err != nil {
						log.Printf("polling error: %v\n", err)
						return
					}
					tx.Reward = reward
					tx.Status = transaction.PROCESSED

					if err := t.repo.Update(ctx, tx); err != nil {
						log.Printf("polling error: %v\n", err)
						return
					}
				}(wg, &tx)
			}
			wg.Wait()

			if err := t.repo.RewardAccounts(ctx, maxUnprocessedTransactionPerQuery); err != nil {
				return err
			}

		}
	}
}

func NewTransactionService(transactions repository.Transaction, cfg *config.AppConfig) *TransactionService {
	return &TransactionService{repo: transactions, cfg: cfg}
}
