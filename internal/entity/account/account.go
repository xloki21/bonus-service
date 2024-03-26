package account

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"math/rand"
)

// UserID UUIDv string user identifier
type UserID string

func (u UserID) Validate() error {
	if _, err := uuid.Parse(string(u)); err != nil {
		return err
	}
	return nil
}

// Account user account in loyalty program
type Account struct {
	ID      UserID `json:"user_id" bson:"user_id"`
	Balance int    `json:"balance" bson:"balance"`
}

func (u *Account) Validate() error {
	// Validate UserID
	if err := u.ID.Validate(); err != nil {
		return err
	}
	// Validate Balance
	if u.Balance < 0 {
		return apperr.AccountInvalidBalance
	}

	return nil
}

func TestAccount() *Account {
	return &Account{
		ID:      UserID(uuid.NewString()),
		Balance: rand.Intn(2104),
	}
}
