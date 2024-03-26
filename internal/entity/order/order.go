package order

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"time"
)

const MaxOrderGoodsAmount = 100

type GoodID string

func (id GoodID) Validate() error {
	if _, err := uuid.Parse(string(id)); err != nil {
		return err
	}
	return nil
}

// Order Пользовательский заказ-покупка
// - все поля обязтельны
// - элементы goods уникалены, даже если пользователь купил два товара, начисление производится только на один
// - масмимальное количество элементов goods - 100
// - время всегда в UTC
type Order struct {
	UserID    account.UserID `json:"user_id" bson:"user_id"`
	Goods     []GoodID       `json:"goods" bson:"goods"`
	Timestamp int64          `json:"timestamp" bson:"timestamp"`
}

func (o Order) Validate() error {
	if err := o.UserID.Validate(); err != nil {
		return apperr.OrderValidationFailed
	}

	if len(o.Goods) > MaxOrderGoodsAmount || len(o.Goods) == 0 {
		return apperr.OrderValidationFailed
	}

	uniques := make(map[GoodID]bool)
	for _, good := range o.Goods {
		if err := good.Validate(); err != nil {
			return apperr.OrderValidationFailed
		}

		if _, exists := uniques[good]; exists {
			return apperr.OrderValidationFailed
		}
		uniques[good] = true
	}
	return nil
}

func TestOrder(goodsAmount int) *Order {
	goods := make([]GoodID, 0, goodsAmount)
	for i := 0; i < goodsAmount; i++ {
		goods = append(goods, GoodID(uuid.NewString()))
	}
	return &Order{
		UserID:    account.UserID(uuid.NewString()),
		Goods:     goods,
		Timestamp: time.Now().Unix(),
	}
}
