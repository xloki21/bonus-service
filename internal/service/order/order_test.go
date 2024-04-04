package order

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/pkg/log"
	"testing"
)

func TestNewOrderService_Register(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("register New order", func(t *testing.T) {
		t.Parallel()
		mock := NewMockorderRepository(ctrl)
		s := NewOrderService(mock)

		testOrder := order.TestOrder(100)
		mock.EXPECT().Register(gomock.Any(), testOrder).Return(nil)
		assert.Nil(t, s.Register(ctx, testOrder), "Should be no error")
	})

	t.Run("register Already registered order", func(t *testing.T) {
		t.Parallel()
		mock := NewMockorderRepository(ctrl)
		s := NewOrderService(mock)

		testOrder := order.TestOrder(100)
		mock.EXPECT().Register(gomock.Any(), testOrder).Return(nil)
		_ = s.Register(ctx, testOrder)

		mock.EXPECT().Register(gomock.Any(), testOrder).Return(apperr.OrderAlreadyRegistered)
		assert.ErrorIs(t, s.Register(ctx, testOrder), apperr.OrderAlreadyRegistered)
	})

}
