package transaction

const (
	COMPLETED   = "COMPLETED"
	PROCESSED   = "PROCESSED"
	UNPROCESSED = "UNPROCESSED"
)

// AggregatedTransaction struct is used to operate with grouped transactions.
type AggregatedTransaction struct {
	ID           interface{}   `bson:"_id,omitempty"`
	UserID       string        `bson:"user_id"`
	Reward       uint          `bson:"reward"`
	Transactions []interface{} `bson:"transactions"`
}

// DTO struct is used to operate with transactions in storage.
type DTO struct {
	ID           interface{} `bson:"_id,omitempty"`
	UserID       string      `bson:"user_id"`
	GoodID       string      `bson:"good_id"`
	Status       string      `bson:"status"`
	Timestamp    int64       `bson:"timestamp"`
	RegisteredAt int64       `bson:"registered_at"`
	CompletedAt  int64       `bson:"completed_at,omitempty"`
	Reward       uint        `bson:"reward,omitempty"`
}

// Transaction struct is used to represent transaction data.
type Transaction struct {
	ID           interface{}
	UserID       string
	GoodID       string
	Status       string
	Timestamp    int64
	RegisteredAt int64
	CompletedAt  int64
	Reward       uint
}

// ToDTO converts Transaction to DTO.
func (t Transaction) ToDTO() DTO {
	return DTO{
		ID:           t.ID,
		UserID:       t.UserID,
		GoodID:       t.GoodID,
		Status:       t.Status,
		Timestamp:    t.Timestamp,
		RegisteredAt: t.RegisteredAt,
		CompletedAt:  t.CompletedAt,
		Reward:       t.Reward,
	}
}
