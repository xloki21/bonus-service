package service

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/repository/mongodb"
	"github.com/xloki21/bonus-service/pkg/log"
	"testing"
)

func TestAccountService_Credit(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := mongodb.NewMongoDB(context.Background(), mongodb.TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			panic(err)
		}
	}()

	repo := mongodb.NewAccountMongoDB(db)
	s := NewAccountService(repo, log.GetDefaultLogger(&config.LoggerConfig{
		Level:    "error",
		Encoding: "console",
	}))

	type args struct {
		id    account.UserID
		value int
	}

	type testCase struct {
		name        string
		args        args
		expectedErr error
	}

	testAccount, err := s.CreateAccount(ctx, 0)
	if err != nil {
		t.Fatalf("failed to create test account: %v", err)
	}

	testCases := []testCase{
		{
			name:        "with positive value",
			args:        args{id: testAccount.ID, value: 100},
			expectedErr: nil,
		},
		{
			name:        "with negative value",
			args:        args{id: testAccount.ID, value: -100},
			expectedErr: apperr.InvalidCreditValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := s.Credit(ctx, tc.args.id, tc.args.value); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestAccountService_Debit(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := mongodb.NewMongoDB(context.Background(), mongodb.TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			panic(err)
		}
	}()

	repo := mongodb.NewAccountMongoDB(db)
	s := NewAccountService(repo, log.GetDefaultLogger(&config.LoggerConfig{
		Level:    "error",
		Encoding: "console",
	}))

	type args struct {
		id    account.UserID
		value int
	}

	type testCase struct {
		name        string
		args        args
		expectedErr error
	}

	testAccountBalance := 21000
	testAccount, err := s.CreateAccount(ctx, testAccountBalance)
	if err != nil {
		t.Fatalf("failed to create test account: %v", err)
	}

	testCases := []testCase{
		{
			name:        "with negative value",
			args:        args{id: testAccount.ID, value: -21},
			expectedErr: apperr.InvalidDebitValue,
		},
		{
			name:        "insufficient balance case",
			args:        args{id: testAccount.ID, value: testAccountBalance + 1},
			expectedErr: apperr.InsufficientBalance,
		},
		{
			name:        "sufficient balance case",
			args:        args{id: testAccount.ID, value: testAccountBalance - 1},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := s.Debit(ctx, tc.args.id, tc.args.value); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestAccountService_CreateAccount(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := mongodb.NewMongoDB(context.Background(), mongodb.TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			panic(err)
		}
	}()

	repo := mongodb.NewAccountMongoDB(db)
	s := NewAccountService(repo, log.GetDefaultLogger(&config.LoggerConfig{
		Level:    "error",
		Encoding: "console",
	}))

	type testCase struct {
		name        string
		value       int
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "with positive balance",
			value:       100,
			expectedErr: nil,
		},
		{
			name:        "with negative balance",
			value:       -100,
			expectedErr: apperr.AccountInvalidBalance,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := s.CreateAccount(ctx, tc.value); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
