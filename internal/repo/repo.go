package repo

import (
	"context"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/repo/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination=mocks/mocks.go -source=repo.go -package=mocks
type Order interface {
	Register(ctx context.Context, o order.Order) error
}

type Account interface {
	Create(context.Context, account.Account) error
	Delete(context.Context, account.Account) error
	FindByID(context.Context, account.UserID) (*account.Account, error)
	Credit(context.Context, account.UserID, uint) error
	Debit(context.Context, account.UserID, uint) error
}

type Transaction interface {
	FindUnprocessed(ctx context.Context, limit int64) ([]transaction.Transaction, error)
	GetOrderTransactions(context.Context, order.Order) ([]transaction.Transaction, error)
	RewardAccounts(ctx context.Context, limit int64) error
	Update(ctx context.Context, tx *transaction.Transaction) error
}

type Repository struct {
	Account
	Order
	Transaction
}

func NewRepositoryMongoDB(db *mongo.Database) *Repository {
	return &Repository{
		Account:     mongodb.NewAccountMongoDB(db),
		Order:       mongodb.NewOrderMongoDB(db),
		Transaction: mongodb.NewTransactionMongoDB(db),
	}
}
