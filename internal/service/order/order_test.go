package order

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/service/order/mock"
	"github.com/xloki21/bonus-service/pkg/log"
	"testing"
)

func TestNewOrderService_Register(t *testing.T) {
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrder := mock.NewMockOrder(ctrl)
	s := NewOrderService(mockOrder)

	testOrder := order.TestOrder(100)

	mockOrder.EXPECT().Register(gomock.Any(), testOrder).Return(nil)
	mockOrder.EXPECT().Register(gomock.Any(), testOrder).Return(apperr.OrderAlreadyRegistered)

	assert.Nil(t, s.Register(ctx, testOrder), "Register New Account: should be no error")
	assert.ErrorIs(t, s.Register(ctx, testOrder), apperr.OrderAlreadyRegistered, "should be an error")
}
