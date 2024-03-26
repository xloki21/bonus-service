package service

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/repository"
)

type OrderService struct {
	orders repository.Order
}

func (o *OrderService) Register(ctx context.Context, order *order.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}
	return o.orders.Register(ctx, order)
}

func NewOrderService(orders repository.Order) *OrderService {
	return &OrderService{orders: orders}
}
