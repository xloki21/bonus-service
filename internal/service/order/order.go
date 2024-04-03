package order

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repository"
)

type Order interface {
	Register(context.Context, *order.Order) error
}

type Service struct {
	orders repository.Order
}

func (o *Service) Register(ctx context.Context, order *order.Order) error {
	if err := order.Validate(); err != nil {
		//o.logger.Warnf("order validation failed: %v", err)
		return err
	}

	if err := o.orders.Register(ctx, order); err != nil {
		//o.logger.Warnf("order registration failed: %v", err)
		return err
	}
	//o.logger.Info("order successfully registered")
	return nil
}

func NewOrderService(orders repository.Order) *Service {
	return &Service{orders: orders}
}
