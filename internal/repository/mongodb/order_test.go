package mongodb

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"testing"
)

func TestOrderMongoDB_Register(t *testing.T) {
	ctx := context.Background()
	db, teardown, err := NewMongoDB(context.Background(), TestDBConfig)

	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err := teardown(ctx); err != nil {
			panic(err)
		}
	}()

	or := NewOrderMongoDB(db)
	tr := NewTransactionMongoDB(db)
	type args struct {
		order *order.Order
	}

	type testCase struct {
		name          string
		args          args
		precondition  func() error
		postcondition func() error
		expectedErr   error
	}

	testOrder := order.TestOrder(rand.Intn(1000) + 1)

	testCases := []testCase{
		{
			name:        "new order",
			args:        args{order: order.TestOrder(10)},
			expectedErr: nil,
		},
		{
			name: "already registered order",
			precondition: func() error {
				return or.Register(ctx, testOrder)
			},
			args:        args{order: testOrder},
			expectedErr: apperr.OrderAlreadyRegistered,
		},
	}

	countDocuments := func(ctx context.Context, collectionName string) (int64, error) {
		return db.Collection(collectionName).CountDocuments(ctx, bson.M{})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				if err := tc.precondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
			ordersDocsBefore, err := countDocuments(ctx, ordersCollection)
			if err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}

			txDocsBefore, err := tr.GetOrderTransactions(ctx, tc.args.order)
			if err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}
			if err := or.Register(ctx, tc.args.order); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}

			ordersDocsAfter, err := countDocuments(ctx, ordersCollection)
			if err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}
			txDocsAfter, err := tr.GetOrderTransactions(ctx, tc.args.order)
			if err != nil {
				t.Errorf("expected error %v, got %v", nil, err)
			}

			if tc.expectedErr != nil {
				if ordersDocsBefore != ordersDocsAfter {
					t.Errorf("expected order collection size not changed, but size has changed from %v to %v",
						ordersDocsBefore, ordersDocsAfter)
				}
				if len(txDocsBefore) != len(txDocsAfter) {
					t.Errorf("expected transactions collection size not changed, but size has changed from %v to %v",
						ordersDocsBefore, ordersDocsAfter)
				}
			} else {
				if (ordersDocsAfter - ordersDocsBefore) != 1 {
					t.Errorf("expected order collection docs to be updated with 1 document, but got %v", ordersDocsAfter-ordersDocsBefore)
				}

				if len(txDocsAfter)-len(txDocsBefore) != len(tc.args.order.Goods) {
					t.Errorf("expected transactions collection docs to be updated with %v document, but got %v", len(tc.args.order.Goods), len(txDocsAfter)-len(txDocsBefore))
				}

			}

			if tc.postcondition != nil {
				if err := tc.postcondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}

		})
	}

}
