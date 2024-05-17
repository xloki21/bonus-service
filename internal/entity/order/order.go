package order

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"time"
)

const MaxOrderGoodsAmount = 100

// DTO struct is used to operate with order data in storage
type DTO struct {
	UserID    string   `bson:"user_id"`
	Goods     []string `bson:"goods"`
	Timestamp int64    `bson:"timestamp"`
}

// Order struct is used to represent order data.
type Order struct {
	UserID    string
	Goods     []string
	Timestamp int64
}

// ToDTO converts Order to DTO.
func (o Order) ToDTO() DTO {
	return DTO{
		UserID:    o.UserID,
		Goods:     o.Goods,
		Timestamp: o.Timestamp,
	}
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

func (o Order) GetTransactions() []transaction.DTO {
	txs := make([]transaction.DTO, 0, len(o.Goods))
	registeredAt := time.Now().Unix()
	for _, goodID := range o.Goods {
		tx := transaction.DTO{
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
