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

func TestService_Register(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	log.BuildLogger(log.TestLoggerConfig)
	ctrl := gomock.NewController(t)

	t.Run("register New order", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockOrder(ctrl)
		s := NewOrderService(mock)

		testOrder := faker.NewOrder(100)

		mock.
			EXPECT().
			Register(gomock.Any(), gomock.Eq(testOrder.ToDTO())).
			Return(nil)

		assert.NoError(t, s.Register(ctx, testOrder))
	})

	t.Run("register already registered order", func(t *testing.T) {
		t.Parallel()
		mock := mocks.NewMockOrder(ctrl)
		s := NewOrderService(mock)

		testOrder := faker.NewOrder(100)
		mock.
			EXPECT().
			Register(gomock.Any(), gomock.Eq(testOrder.ToDTO())).
			Return(nil)

		err := s.Register(ctx, testOrder)
		assert.NoError(t, err)

		mock.
			EXPECT().
			Register(gomock.Any(), gomock.Eq(testOrder.ToDTO())).
			Return(apperr.OrderAlreadyRegistered)

		assert.ErrorIs(t, s.Register(ctx, testOrder), apperr.OrderAlreadyRegistered)
	})

}
