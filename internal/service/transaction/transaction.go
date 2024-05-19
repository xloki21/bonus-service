package transaction

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/client/accrual"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/repo"
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
	client *accrual.Client
	repo   repo.Transaction
	cfg    config.TransactionServiceConfig
}

func txProcessingRound(ctx context.Context, repo repo.Transaction, client *accrual.Client, batchSize int64) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}
	logger.Info("polling event triggered")
	logger.Info("find unprocessed transactions...")
	txs, err := repo.FindUnprocessed(ctx, batchSize)
	if err != nil {
		return errors.Join(apperr.TransactionProcessingError, err)
	}
	logger.Info("processing transactions...")

	if len(txs) == 0 {
		return nil
	}
	errsCh := make(chan error, len(txs)) // channel to collect errors

	wg := &sync.WaitGroup{}
	for index := range txs {
		wg.Add(1)
		go func(wg *sync.WaitGroup, index int) {
			defer wg.Done()
			reward, err := client.GetAccrual(ctx, &txs[index])
			if err != nil {
				errsCh <- errors.Join(apperr.TransactionProcessingError, err)
				return
			}

			txs[index].Reward = reward
			txs[index].Status = transaction.PROCESSED

			if err := repo.Update(ctx, &txs[index]); err != nil {
				errj := errors.Join(apperr.TransactionProcessingError, err)
				errsCh <- errj
				return
			}
		}(wg, index)
	}

	wg.Wait()
	close(errsCh)

	for err := range errsCh {
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
	}
	return nil
}

// Polling is a blocking method that polls unprocessed transactions.
func (s *Service) Polling(ctx context.Context) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}
	successfulRounds := 0
	logger.Info("polling transactions...")
	ticker := time.NewTicker(s.cfg.PollingInterval)
	defer ticker.Stop() // Stop the ticker so it can be garbage collected

	batchSize := int64(s.cfg.MaxTransactionsPerRequest)
	for {
		select {
		case <-ctx.Done():
			logger.Info("polling events listener stopped")
			return ctx.Err()
		case <-ticker.C:

			successfulRounds = successfulRounds + 1
			if err := txProcessingRound(ctx, s.repo, s.client, batchSize); err != nil {
				successfulRounds = 0
				if errors.Is(err, apperr.AccrualServiceTooManyRequests) {
					adjustedRPS := int(float32(s.client.GetRPS()) * 0.95)
					s.client.AdjustRPS(adjustedRPS)
					break
				}
			}

			if successfulRounds > minSuccessfulRoundsToRestoreRPS {
				// try to restore RPS after successful rounds
				successfulRounds = minSuccessfulRoundsToRestoreRPS
				adjustedRPS := int(float32(s.client.GetRPS()) * 1.05)
				adjustedRPS = int(math.Min(float64(adjustedRPS), float64(maxRequestsPerSecond))) // cap
				s.client.AdjustRPS(adjustedRPS)
			}

			logger.Info("rewarding accounts...")
			if err := s.repo.RewardAccounts(ctx, batchSize); err != nil {
				logger.Warnf("polling event error on reward accounts: %v", err)
			}

		}
	}
}

func NewTransactionService(transactions repo.Transaction, client *accrual.Client, cfg config.TransactionServiceConfig) *Service {
	return &Service{repo: transactions, client: client, cfg: cfg}
}
