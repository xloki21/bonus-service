package mongodb

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/faker"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestTransactionMongoDB_FindUnprocessed(t *testing.T) {
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

	or := NewOrderStorage(db)
	tr := NewTransactionStorage(db)
	type args struct {
		order order.Order
		limit int64
	}

	type testCase struct {
		name          string
		args          args
		precondition  func() error
		postcondition func() error
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "new order transactions: len(tx) > limit",
			args: args{order: faker.NewOrder(1000), limit: 10},
			postcondition: func() error {
				_, err := or.db.Collection(ordersCollection).DeleteMany(ctx, bson.M{})
				return err
			},
			expectedErr: nil,
		},
		{
			name: "new order transactions: len(tx) < limit",
			args: args{order: faker.NewOrder(1000), limit: 2000},
			postcondition: func() error {
				_, err := or.db.Collection(ordersCollection).DeleteMany(ctx, bson.M{})
				return err
			},
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

			txBefore, err := tr.FindUnprocessed(ctx, tc.args.limit)
			if err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}

			if err := or.Register(ctx, tc.args.order); err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}
			txAfter, err := tr.FindUnprocessed(ctx, tc.args.limit)
			if err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}
			newTransactions := len(txAfter) - len(txBefore)
			newTxExpected := min(int(tc.args.limit), newTransactions)
			if newTransactions != newTxExpected {
				t.Errorf("expected new transactions %v, got %v", len(tc.args.order.Goods), newTxExpected)
			}

			if tc.postcondition != nil {
				if err := tc.postcondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
		})
	}
}
