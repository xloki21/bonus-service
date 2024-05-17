package account

// Account struct is used to represent account data.
type Account struct {
	ID      string
	Balance uint
}

// DTO struct is used to operate with account data in storage
type DTO struct {
	ID      string `bson:"user_id"`
	Balance uint   `bson:"balance"`
}

func (a *Account) ToDTO() DTO {
	return DTO{
		ID:      a.ID,
		Balance: a.Balance,
	}
}
