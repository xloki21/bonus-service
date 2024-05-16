package order

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"time"
)

const MaxOrderGoodsAmount = 100

// Order struct is used to store order data.
type Order struct {
	UserID    string   `json:"user_id" bson:"user_id" binding:"required"`
	Goods     []string `json:"goods" bson:"goods" binding:"required"`
	Timestamp int64    `json:"timestamp" bson:"timestamp" binding:"required"`
}

func (o Order) Validate() error {
	if err := uuid.Validate(o.UserID); err != nil {
		return apperr.OrderValidationFailed
	}

	if len(o.Goods) > MaxOrderGoodsAmount || len(o.Goods) == 0 {
		return apperr.OrderValidationFailed
	}

	uniques := make(map[string]bool)
	for _, good := range o.Goods {
		if err := uuid.Validate(good); err != nil {
			return apperr.OrderValidationFailed
		}

		if _, exists := uniques[good]; exists {
			return apperr.OrderValidationFailed
		}
		uniques[good] = true
	}
	return nil
}

func (o Order) GetTransactions() []transaction.Transaction {
	txs := make([]transaction.Transaction, 0, len(o.Goods))
	registeredAt := time.Now().Unix()
	for _, goodID := range o.Goods {
		tx := transaction.Transaction{
			UserID:       o.UserID,
			Status:       transaction.UNPROCESSED,
			GoodID:       goodID,
			Timestamp:    o.Timestamp,
			RegisteredAt: registeredAt,
		}
		txs = append(txs, tx)
	}
	return txs
}
