package account

import (
	"github.com/google/uuid"
)

// Account of user in loyalty program.
type Account struct {
	ID      string `json:"user_id" bson:"user_id"`
	Balance uint   `json:"balance" bson:"balance"`
}

func (u *Account) Validate() error {
	// Validate UserID
	if err := uuid.Validate(u.ID); err != nil {
		return err
	}

	return nil
}
