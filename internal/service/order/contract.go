package order

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/order"
)

//go:generate mockgen -source=contract.go -destination=mock/mock.go -package=mock
type Order interface {
	Register(context.Context, *order.Order) error
}
