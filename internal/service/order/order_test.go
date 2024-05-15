package order

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

func TestNewOrderService_Register(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("register New order", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockOrder(ctrl)
		s := NewOrderService(mock)

		testOrder := faker.NewOrder(100)

		mock.
			EXPECT().
			Register(gomock.Any(), gomock.Eq(testOrder)).
			Return(nil)

		assert.Nil(t, s.Register(ctx, testOrder), "Should be no error")
	})

	t.Run("register Already registered order", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockOrder(ctrl)
		s := NewOrderService(mock)

		testOrder := faker.NewOrder(100)
		mock.
			EXPECT().
			Register(gomock.Any(), gomock.Eq(testOrder)).
			Return(nil)

		err := s.Register(ctx, testOrder)
		assert.NoError(t, err)

		mock.
			EXPECT().
			Register(gomock.Any(), gomock.Eq(testOrder)).
			Return(apperr.OrderAlreadyRegistered)

		assert.ErrorIs(t, s.Register(ctx, testOrder), apperr.OrderAlreadyRegistered)
	})

}
