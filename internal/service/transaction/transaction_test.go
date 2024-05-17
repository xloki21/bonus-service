package transaction

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/client/accrual"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/faker"
	"github.com/xloki21/bonus-service/internal/repo/mocks"
	httpcMocks "github.com/xloki21/bonus-service/pkg/httppc/mocks"
	"github.com/xloki21/bonus-service/pkg/log"
	"io"
	"net/http"
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

		mockHttppc := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(mockHttppc)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)
		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(make([]transaction.DTO, 0, batchSize), nil).AnyTimes()

		err = s.Polling(ctxd)
		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("Successfully processed transactions", func(t *testing.T) {
		t.Parallel()
		ctxd, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)

		httpcMock := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(httpcMock)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)
		testOrderTxs := faker.NewOrder(int(batchSize)).GetTransactions()

		// func call sequence:
		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(testOrderTxs, nil).AnyTimes()

		httpcMock.
			EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(bytes.NewBuffer([]byte("40"))),
				ContentLength: 2,
			}, nil).AnyTimes()

		for i := range testOrderTxs {
			testOrderTxs[i].Reward = 40
			testOrderTxs[i].Status = transaction.PROCESSED

		}

		for i := range testOrderTxs {
			mock.
				EXPECT().
				Update(gomock.Any(), gomock.Eq(&testOrderTxs[i])).
				Return(nil).AnyTimes()
		}

		mock.
			EXPECT().
			RewardAccounts(gomock.Any(), gomock.Eq(batchSize)).
			Return(nil).AnyTimes()

		err = s.Polling(ctxd)
		assert.Equal(t, context.DeadlineExceeded, err)

	})

}
