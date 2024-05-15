package faker

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"time"
)

func NewOrder(goodsAmount int) order.Order {
	goods := make([]string, 0, goodsAmount)
	for i := 0; i < goodsAmount; i++ {
		goods = append(goods, uuid.NewString())
	}
	return order.Order{
		UserID:    uuid.NewString(),
		Goods:     goods,
		Timestamp: time.Now().Unix(),
	}
}
