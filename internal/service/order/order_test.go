package order

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repository/mongodb"
	"testing"
	"time"
)

func TestOrderService_Register(t *testing.T) {
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

	repo := mongodb.NewOrderMongoDB(db)
	s := NewOrderService(repo)
	type args struct {
		order *order.Order
	}

	type testCase struct {
		name         string
		args         args
		precondition func() error
		expectedErr  error
	}

	testOrder := order.TestOrder(10)

	testCases := []testCase{
		{
			name:        "new order with 3 goods",
			args:        args{order: order.TestOrder(3)},
			expectedErr: nil,
		},
		{
			name: "new order with empty goods list",
			args: args{order: &order.Order{
				UserID:    account.UserID(uuid.NewString()),
				Goods:     make([]order.GoodID, 0),
				Timestamp: time.Now().Unix(),
			}},
			expectedErr: apperr.OrderValidationFailed,
		},
		{
			name: "new order with non-unique good indices in list",
			args: args{order: &order.Order{
				UserID: account.UserID(uuid.NewString()),
				Goods: []order.GoodID{
					order.GoodID("81755586-1269-4fe7-8141-6580290767da"),
					order.GoodID("7c552789-9019-41bd-8c79-e73145757445"),
					order.GoodID("81755586-1269-4fe7-8141-6580290767da"),
				},
				Timestamp: time.Now().Unix(),
			}},
			expectedErr: apperr.OrderValidationFailed,
		},
		{
			name: "new order with invalid good indices in list",
			args: args{order: &order.Order{
				UserID: account.UserID(uuid.NewString()),
				Goods: []order.GoodID{
					order.GoodID("81755586-1269-4fe7-8141-6580290767da"),
					order.GoodID("7c552789-9019-41bd-8c79-e73"),
					order.GoodID("81755586-1269-4fe7-8141-6580290767da"),
				},
				Timestamp: time.Now().Unix(),
			}},
			expectedErr: apperr.OrderValidationFailed,
		},
		{
			name:        "new order with exceeding amount of goods",
			args:        args{order: order.TestOrder(order.MaxOrderGoodsAmount + 1)},
			expectedErr: apperr.OrderValidationFailed,
		},
		{
			name: "already registered order",
			args: args{order: testOrder},
			precondition: func() error {
				return s.Register(ctx, testOrder)
			},
			expectedErr: apperr.OrderAlreadyRegistered,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				if err := tc.precondition(); err != nil {
					t.Errorf("expected error %v, got %v", nil, err)
				}
			}
			if err := s.Register(ctx, tc.args.order); !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
