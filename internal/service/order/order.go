package order

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repo"
	"github.com/xloki21/bonus-service/pkg/log"
)

type Service struct {
	orders repo.Order
}

func (o *Service) Register(ctx context.Context, order order.Order) error {
	logger, err := log.GetLogger()
	if err != nil {
		return err
	}

	if err := o.orders.Register(ctx, order.ToDTO()); err != nil {
		logger.Warnf("order registration failed: %v", err)
		return err
	}
	logger.Info("order successfully registered")
	return nil
}

func NewOrderService(orders repo.Order) *Service {
	return &Service{orders: orders}
}
