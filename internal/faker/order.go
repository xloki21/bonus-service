package faker

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"time"
)

func NewOrder(goodsAmount int) order.Order {
	goods := make([]order.GoodID, 0, goodsAmount)
	for i := 0; i < goodsAmount; i++ {
		goods = append(goods, order.GoodID(uuid.NewString()))
	}
	return order.Order{
		UserID:    account.UserID(uuid.NewString()),
		Goods:     goods,
		Timestamp: time.Now().Unix(),
	}
}
