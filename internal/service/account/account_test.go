package account

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/faker"
	"github.com/xloki21/bonus-service/internal/repo/mocks"
	"github.com/xloki21/bonus-service/pkg/log"
	"testing"
)

func TestAccountService_Debit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("debit account with insufficient funds", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount)).
			Return(nil)

		err := s.CreateAccount(ctx, testAccount)
		assert.NoError(t, err)

		value := testAccount.Balance + 1
		mock.
			EXPECT().
			Debit(gomock.Any(), gomock.Eq(testAccount.ID), gomock.Eq(value)).
			Return(apperr.InsufficientBalance)

		assert.ErrorIs(t, s.Debit(ctx, testAccount.ID, value), apperr.InsufficientBalance)
	})

	t.Run("debit account with success", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount)).
			Return(nil)

		err := s.CreateAccount(ctx, testAccount)
		assert.NoError(t, err)

		value := testAccount.Balance

		mock.
			EXPECT().
			Debit(gomock.Any(), gomock.Eq(testAccount.ID), gomock.Eq(value)).
			Return(nil)

		assert.Nil(t, s.Debit(ctx, testAccount.ID, value), "should be no error")
	})
}

func TestAccountService_Credit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("credit account with success", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount)).
			Return(nil)

		err := s.CreateAccount(ctx, testAccount)
		assert.NoError(t, err)

		value := testAccount.Balance

		mock.EXPECT().
			Credit(gomock.Any(), gomock.Eq(testAccount.ID), gomock.Eq(value)).
			Return(nil)

		assert.Nil(t, s.Credit(ctx, testAccount.ID, value), "should be no error")
	})
}

func TestAccountService_CreateAccount(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("create new account with success", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount)).
			Return(nil)

		assert.Nil(t, s.CreateAccount(ctx, testAccount), "should be no error")
	})

	t.Run("create already registered account with fail", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount)).
			Return(nil)

		_ = s.CreateAccount(ctx, testAccount)
		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount)).
			Return(apperr.AccountAlreadyExists)

		assert.ErrorIs(t, s.CreateAccount(ctx, testAccount), apperr.AccountAlreadyExists)
	})
}
