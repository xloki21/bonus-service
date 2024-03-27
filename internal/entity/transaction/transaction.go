package transaction

import (
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/entity/order"
)

type Status string

const (
	COMPLETED   Status = "COMPLETED"
	PROCESSED   Status = "PROCESSED"
	UNPROCESSED Status = "UNPROCESSED"
)

type AggregatedTransaction struct {
	ID           interface{}    `bson:"_id,omitempty"`
	UserID       account.UserID `bson:"user_id"`
	Reward       uint           `bson:"reward"`
	Transactions []interface{}  `bson:"transactions"`
}

type Transaction struct {
	ID           interface{}    `bson:"_id,omitempty"`
	UserID       account.UserID `bson:"user_id"`
	GoodID       order.GoodID   `bson:"good_id"`
	Status       Status         `bson:"status"`
	Timestamp    int64          `bson:"timestamp"`
	RegisteredAt int64          `bson:"registered_at"`
	CompletedAt  int64          `bson:"completed_at,omitempty"`
	Reward       uint           `bson:"reward,omitempty"`
}
