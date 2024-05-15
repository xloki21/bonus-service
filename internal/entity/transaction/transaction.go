package transaction

const (
	COMPLETED   = "COMPLETED"
	PROCESSED   = "PROCESSED"
	UNPROCESSED = "UNPROCESSED"
)

type AggregatedTransaction struct {
	ID           interface{}   `bson:"_id,omitempty"`
	UserID       string        `bson:"user_id"`
	Reward       uint          `bson:"reward"`
	Transactions []interface{} `bson:"transactions"`
}

type Transaction struct {
	ID           interface{} `bson:"_id,omitempty"`
	UserID       string      `bson:"user_id"`
	GoodID       string      `bson:"good_id"`
	Status       string      `bson:"status"`
	Timestamp    int64       `bson:"timestamp"`
	RegisteredAt int64       `bson:"registered_at"`
	CompletedAt  int64       `bson:"completed_at,omitempty"`
	Reward       uint        `bson:"reward,omitempty"`
}
