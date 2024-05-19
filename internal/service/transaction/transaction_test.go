package transaction

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/client/accrual"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/faker"
	"github.com/xloki21/bonus-service/internal/repo/mocks"
	httpcMocks "github.com/xloki21/bonus-service/pkg/httppc/mocks"
	"github.com/xloki21/bonus-service/pkg/log"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestService_Polling(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cfg, err := config.InitConfigFromViper()
	assert.NoError(t, err)
	log.BuildLogger(log.TestLoggerConfig)
	batchSize := int64(cfg.TransactionServiceConfig.MaxTransactionsPerRequest)

	t.Run("unprocessed transactions error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t) // New in go1.14+ no longer need to call ctrl.Finish()
		ctxd, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)

		mockHttppc := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(mockHttppc)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(nil, errors.Join(apperr.TransactionProcessingError, errors.New("find unprocessed transactions error")))

		assert.ErrorIs(t, txProcessingRound(ctxd, s.repo, s.client, batchSize), apperr.TransactionProcessingError)
	})

	t.Run("accrual service error: accrual not found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t) // New in go1.14+ no longer need to call ctrl.Finish()
		ctxd, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)

		mockRequest := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(mockRequest)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		orderTransactions := faker.NewOrder(int(batchSize)).GetTransactions()

		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(orderTransactions, nil)

		mockRequest.
			EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				StatusCode: http.StatusNotFound,
			}, apperr.AccrualNotFound).AnyTimes()

		assert.ErrorIs(t, txProcessingRound(ctxd, s.repo, s.client, batchSize), apperr.TransactionProcessingError)
	})

	t.Run("accrual service error: too many requests", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t) // New in go1.14+ no longer need to call ctrl.Finish()
		ctxd, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)

		mockRequest := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(mockRequest)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		orderTransactions := faker.NewOrder(int(batchSize)).GetTransactions()

		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(orderTransactions, nil)

		mockRequest.
			EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				StatusCode: http.StatusTooManyRequests,
			}, apperr.AccrualServiceTooManyRequests).AnyTimes()

		assert.ErrorIs(t, txProcessingRound(ctxd, s.repo, s.client, batchSize), apperr.TransactionProcessingError)
	})

	t.Run("Success path: no new transactions to process", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t) // New in go1.14+ no longer need to call ctrl.Finish()
		ctxd, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)

		mockRequest := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(mockRequest)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(make([]transaction.DTO, 0, batchSize), nil)

		assert.NoError(t, txProcessingRound(ctxd, s.repo, s.client, batchSize))
	})

	t.Run("Success path: all ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t) // New in go1.14+ no longer need to call ctrl.Finish()

		ctxd, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		mock := mocks.NewMockTransaction(ctrl)
		client := accrual.NewClient(cfg.AccrualService)

		mockRequest := httpcMocks.NewMockHTTPDoer(ctrl)
		client.SetHTTPClient(mockRequest)
		s := NewTransactionService(mock, client, cfg.TransactionServiceConfig)

		beforeTestOrderTxs := faker.NewOrder(int(batchSize)).GetTransactions()
		afterTestOrderTxs := make([]transaction.DTO, batchSize)

		copy(afterTestOrderTxs, beforeTestOrderTxs)
		mock.
			EXPECT().
			FindUnprocessed(gomock.Any(), gomock.Eq(batchSize)).
			Return(beforeTestOrderTxs, nil)

		testRewardValue := uint(40)
		for _ = range beforeTestOrderTxs {
			mockRequest.
				EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode:    http.StatusOK,
					Body:          io.NopCloser(bytes.NewBuffer([]byte(strconv.Itoa(int(testRewardValue))))),
					ContentLength: 2,
				}, nil)
		}

		for i := range afterTestOrderTxs {
			afterTestOrderTxs[i].Reward = testRewardValue
			afterTestOrderTxs[i].Status = transaction.PROCESSED
		}

		for i := range afterTestOrderTxs {
			mock.
				EXPECT().
				Update(gomock.Any(), gomock.Eq(&afterTestOrderTxs[i])).
				Return(nil)
		}

		assert.NoError(t, txProcessingRound(ctxd, s.repo, s.client, batchSize))
	})
}
