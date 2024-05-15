package mongodb

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/faker"
	"testing"
)

func TestAccountMongoDB_Create(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	r := NewAccountStorage(db)

	type args struct {
		account account.Account
	}

	type testCase struct {
		name         string
		args         args
		precondition func() error
		expectedErr  error
	}

	testAccount := faker.NewAccount()

	testCases := []testCase{
		{
			name:        "regular account",
			args:        args{account: faker.NewAccount()},
			expectedErr: nil,
		},
		{
			name: "account with registered user id",
			precondition: func() error {
				return r.Create(ctx, testAccount)
			},
			args:        args{account: testAccount},
			expectedErr: apperr.AccountAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				if err := tc.precondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
			if err := r.Create(ctx, tc.args.account); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}

}

func TestAccountMongoDB_FindByID(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	r := NewAccountStorage(db)
	type args struct {
		id account.UserID
	}

	type testCase struct {
		name         string
		args         args
		precondition func() error
		expectedErr  error
	}

	testAccount := faker.NewAccount()

	testCases := []testCase{
		{
			name:        "missing account",
			args:        args{id: account.UserID(uuid.NewString())},
			expectedErr: apperr.AccountNotFound,
		},
		{
			name: "registered account",
			precondition: func() error {
				return r.Create(ctx, testAccount)
			},
			args:        args{id: testAccount.ID},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				if err := tc.precondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
			if _, err := r.FindByID(ctx, tc.args.id); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestAccountMongoDB_Credit(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	r := NewAccountStorage(db)

	type args struct {
		id    account.UserID
		value uint
	}

	type testCase struct {
		name         string
		args         args
		precondition func() error
		expectedErr  error
	}

	testAccount := faker.NewAccount()

	testCases := []testCase{
		{
			name:        "unknown account id",
			args:        args{id: account.UserID(uuid.NewString()), value: 100},
			expectedErr: apperr.AccountNotFound,
		},
		{
			name: "existing account id",
			precondition: func() error {
				return r.Create(ctx, testAccount)
			},
			args:        args{id: testAccount.ID, value: 100},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				if err := tc.precondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
			if err := r.Credit(ctx, tc.args.id, tc.args.value); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestAccountMongoDB_Debit(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	r := NewAccountStorage(db)

	type args struct {
		id    account.UserID
		value uint
	}

	type testCase struct {
		name          string
		args          args
		precondition  func() error
		postcondition func() error
		expectedErr   error
	}

	testAccount := faker.NewAccount()

	testCases := []testCase{
		{
			name:        "unknown account id",
			args:        args{id: account.UserID(uuid.NewString()), value: 100},
			expectedErr: apperr.InsufficientBalance,
		},
		{
			name: "existing account with insufficient balance",
			precondition: func() error {
				return r.Create(ctx, testAccount)
			},
			postcondition: func() error {
				return r.Delete(ctx, testAccount)
			},
			args:        args{id: testAccount.ID, value: testAccount.Balance + 1},
			expectedErr: apperr.InsufficientBalance,
		},
		{
			name: "existing account with sufficient balance",
			precondition: func() error {
				return r.Create(ctx, testAccount)
			},
			postcondition: func() error {
				return r.Delete(ctx, testAccount)
			},
			args:        args{id: testAccount.ID, value: testAccount.Balance - 1},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				if err := tc.precondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
			if err := r.Debit(ctx, tc.args.id, tc.args.value); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.postcondition != nil {
				if err := tc.postcondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}

		})
	}
}
