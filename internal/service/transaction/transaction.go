package transaction

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/integration"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/pkg/log"
	"math"
	"sync"
	"time"
)

const (
	minSuccessfulRoundsToRestoreRPS = 10
	maxRequestsPerSecond            = 20
)

type Transaction interface {
	Polling(ctx context.Context) error
}

type Service struct {
	cfg  *config.AppConfig
	repo repository.Transaction
}

// Polling is a blocking method that polls unprocessed transactions.
func (t *Service) Polling(ctx context.Context) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}
	successfulRounds := 0
	logger.Info("polling transactions...")
	accrualServiceClient := integration.New(t.cfg.AccrualService)
	ticker := time.NewTicker(t.cfg.TransactionServiceConfig.PollingInterval)
	defer ticker.Stop() // Stop the ticker so it can be garbage collected
	for {
		select {
		case <-ctx.Done():
			logger.Info("polling events listener stopped")
			return ctx.Err()
		case <-ticker.C:
			logger.Info("polling event triggered")
			logger.Info("find unprocessed transactions...")
			txs, err := t.repo.FindUnprocessed(ctx, int64(t.cfg.TransactionServiceConfig.MaxTransactionsPerRequest))
			if err != nil {
				logger.Warnf("polling event error on find unprocessed transactions: %v", err)
				continue
			}
			logger.Info("processing transactions...")

			if len(txs) == 0 {
				continue
			}
			errsCh := make(chan error, len(txs)) // channel to collect errors

			wg := &sync.WaitGroup{}
			for index := range txs {
				wg.Add(1)
				go func(wg *sync.WaitGroup, index int) {
					defer wg.Done()
					reward, err := accrualServiceClient.GetAccrual(ctx, &txs[index])
					if err != nil {
						errsCh <- err
						logger.Warnf("error during request to accrual service: %v", err)
						return
					}

					txs[index].Reward = reward
					txs[index].Status = transaction.PROCESSED

					if err := t.repo.Update(ctx, &txs[index]); err != nil {
						errsCh <- err
						logger.Warnf("polling event error on update transaction: %v", err)
						return
					}
				}(wg, index)
			}
			wg.Wait()
			close(errsCh)

			// if any `AccrualServiceTooManyRequests` occurs, adjust RPS to avoid too many requests
			successfulRounds = successfulRounds + 1
			for err := range errsCh {
				if err != nil {
					successfulRounds = 0
					if errors.Is(err, apperr.AccrualServiceTooManyRequests) {
						adjustedRPS := int(float32(accrualServiceClient.Config.RPS) * 0.95)
						accrualServiceClient.AdjustRPS(adjustedRPS)
						break
					}
				}
			}

			if successfulRounds > minSuccessfulRoundsToRestoreRPS {
				// try to restore RPS after successful rounds
				successfulRounds = minSuccessfulRoundsToRestoreRPS
				adjustedRPS := int(float32(accrualServiceClient.Config.RPS) * 1.05)
				adjustedRPS = int(math.Min(float64(adjustedRPS), float64(maxRequestsPerSecond))) // cap
				accrualServiceClient.AdjustRPS(adjustedRPS)
			}

			logger.Info("rewarding accounts...")
			if err := t.repo.RewardAccounts(ctx, int64(t.cfg.TransactionServiceConfig.MaxTransactionsPerRequest)); err != nil {
				logger.Warnf("polling event error on reward accounts: %v", err)
			}
		}
	}
}

func NewTransactionService(transactions repository.Transaction, cfg *config.AppConfig) *Service {
	return &Service{repo: transactions, cfg: cfg}
}