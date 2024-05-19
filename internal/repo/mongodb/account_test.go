package mongodb

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/faker"
	"testing"
)

func TestAccountMongoDB_Create(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, teardown(ctx))
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
				return r.Create(ctx, testAccount.ToDTO())
			},
			args:        args{account: testAccount},
			expectedErr: apperr.AccountAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				assert.NoError(t, tc.precondition())
			}
			if err := r.Create(ctx, tc.args.account.ToDTO()); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}

}

func TestAccountMongoDB_FindByID(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, teardown(ctx))
	}()

	r := NewAccountStorage(db)
	type args struct {
		id string
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
			args:        args{id: uuid.NewString()},
			expectedErr: apperr.AccountNotFound,
		},
		{
			name: "registered account",
			precondition: func() error {
				return r.Create(ctx, testAccount.ToDTO())
			},
			args:        args{id: testAccount.ID},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				assert.NoError(t, tc.precondition())
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
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, teardown(ctx))
	}()

	r := NewAccountStorage(db)

	type args struct {
		id    string
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
			args:        args{id: uuid.NewString(), value: 100},
			expectedErr: apperr.AccountNotFound,
		},
		{
			name: "existing account id",
			precondition: func() error {
				return r.Create(ctx, testAccount.ToDTO())
			},
			args:        args{id: testAccount.ID, value: 100},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				assert.NoError(t, tc.precondition())
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
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, teardown(ctx))
	}()

	r := NewAccountStorage(db)

	type args struct {
		id    string
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
			args:        args{id: uuid.NewString(), value: 100},
			expectedErr: apperr.InsufficientBalance,
		},
		{
			name: "existing account with insufficient balance",
			precondition: func() error {
				return r.Create(ctx, testAccount.ToDTO())
			},
			postcondition: func() error {
				return r.Delete(ctx, testAccount.ToDTO())
			},
			args:        args{id: testAccount.ID, value: testAccount.Balance + 1},
			expectedErr: apperr.InsufficientBalance,
		},
		{
			name: "existing account with sufficient balance",
			precondition: func() error {
				return r.Create(ctx, testAccount.ToDTO())
			},
			postcondition: func() error {
				return r.Delete(ctx, testAccount.ToDTO())
			},
			args:        args{id: testAccount.ID, value: testAccount.Balance - 1},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				assert.NoError(t, tc.precondition())
			}
			if err := r.Debit(ctx, tc.args.id, tc.args.value); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.postcondition != nil {
				assert.NoError(t, tc.postcondition())
			}

		})
	}
}
