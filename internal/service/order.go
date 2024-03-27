package service

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/pkg/log"
	"github.com/xloki21/bonus-service/internal/repository"
)

type OrderService struct {
	orders repository.Order
	logger log.Logger
}

func (o *OrderService) Register(ctx context.Context, order *order.Order) error {
	if err := order.Validate(); err != nil {
		o.logger.Warnf("order validation failed: %v", err)
		return err
	}

	if err := o.orders.Register(ctx, order); err != nil {
		o.logger.Warnf("order registration failed: %v", err)
		return err
	}
	o.logger.Info("order successfully registered")
	return nil
}

func NewOrderService(orders repository.Order, logger log.Logger) *OrderService {
	return &OrderService{orders: orders, logger: logger}
}
