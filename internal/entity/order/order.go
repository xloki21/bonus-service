package order

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
)

const MaxOrderGoodsAmount = 100

type GoodID string

// Validate validate GoodID.
func (id GoodID) Validate() error {
	if _, err := uuid.Parse(string(id)); err != nil {
		return err
	}
	return nil
}

// Order struct is used to store order data.
type Order struct {
	UserID    account.UserID `json:"user_id" bson:"user_id" binding:"required"`
	Goods     []GoodID       `json:"goods" bson:"goods" binding:"required"`
	Timestamp int64          `json:"timestamp" bson:"timestamp" binding:"required"`
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
