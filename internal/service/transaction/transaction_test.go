package transaction

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/client/accrual"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/faker"
	"github.com/xloki21/bonus-service/internal/repo/mocks"
	"github.com/xloki21/bonus-service/pkg/log"
	"sync"
	"testing"
	"time"
)

func TestService_Polling(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cfg, err := config.InitConfigFromViper()
	assert.NoError(t, err)
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	batchSize := int64(cfg.TransactionServiceConfig.MaxTransactionsPerRequest)
	t.Run("No new transactions to process", func(t *testing.T) {
		t.Parallel()
		ctxd, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(make([]transaction.Transaction, 0, batchSize), nil).AnyTimes()

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Polling(ctxd)
			assert.Equal(t, context.DeadlineExceeded, err)
		}()
		wg.Wait()
	})

	t.Run("Successfully processed transactions", func(t *testing.T) {
		t.Parallel()
		ctxd, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		testOrderTxs := faker.NewOrder(int(batchSize)).GetTransactions()

		findUnprocessed := mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(testOrderTxs, nil).AnyTimes()

		updateTransactions := mock.
			EXPECT().
			Update(gomock.Any(), gomock.Eq(testOrderTxs)).
			Return(nil).AnyTimes()

		rewardAccounts := mock.
			EXPECT().
			RewardAccounts(gomock.Any(), gomock.Eq(batchSize)).
			Return(nil).AnyTimes()

		gomock.InOrder(findUnprocessed, updateTransactions, rewardAccounts)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = s.Polling(ctxd)
			assert.Equal(t, context.DeadlineExceeded, err)
		}()
		wg.Wait()

	})

}
