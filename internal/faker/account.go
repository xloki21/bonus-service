package faker

import (
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"math/rand"
)

func NewAccount() account.Account {
	return account.Account{
		ID:      uuid.NewString(),
		Balance: uint(rand.Intn(2104)),
	}
}
