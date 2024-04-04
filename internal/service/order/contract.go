package order

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=order
type orderRepository interface {
	Register(context.Context, order.Order) error
}
