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

func TestService_Debit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)

	t.Run("debit account with insufficient funds", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount.ToDTO())).
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
			Create(gomock.Any(), gomock.Eq(testAccount.ToDTO())).
			Return(nil)

		err := s.CreateAccount(ctx, testAccount)
		assert.NoError(t, err)

		value := testAccount.Balance

		mock.
			EXPECT().
			Debit(gomock.Any(), gomock.Eq(testAccount.ID), gomock.Eq(value)).
			Return(nil)

		assert.NoError(t, s.Debit(ctx, testAccount.ID, value))
	})
}

func TestService_Credit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)

	t.Run("credit account with success", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount.ToDTO())).
			Return(nil)

		err := s.CreateAccount(ctx, testAccount)
		assert.NoError(t, err)

		value := testAccount.Balance

		mock.EXPECT().
			Credit(gomock.Any(), gomock.Eq(testAccount.ID), gomock.Eq(value)).
			Return(nil)

		assert.NoError(t, s.Credit(ctx, testAccount.ID, value))
	})
}

func TestService_CreateAccount(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)

	t.Run("create new account with success", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount.ToDTO())).
			Return(nil)

		assert.NoError(t, s.CreateAccount(ctx, testAccount))
	})

	t.Run("create already registered account with fail", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockAccount(ctrl)
		s := NewAccountService(mock)
		testAccount := faker.NewAccount()

		mock.EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount.ToDTO())).
			Return(nil)

		_ = s.CreateAccount(ctx, testAccount)
		mock.
			EXPECT().
			Create(gomock.Any(), gomock.Eq(testAccount.ToDTO())).
			Return(apperr.AccountAlreadyExists)

		assert.ErrorIs(t, s.CreateAccount(ctx, testAccount), apperr.AccountAlreadyExists)
	})
}
