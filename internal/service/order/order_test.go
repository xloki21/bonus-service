package order

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repository/mongodb"
	"github.com/xloki21/bonus-service/pkg/log"
	"testing"
)

func TestOrderService_Register(t *testing.T) {
	log.BuildLogger(&log.TestLoggerConfig)

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
