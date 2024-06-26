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
	Register(ctx context.Context, o order.DTO) error
}

type Account interface {
	Create(context.Context, account.DTO) error
	Delete(context.Context, account.DTO) error
	FindByID(context.Context, string) (*account.DTO, error)
	Credit(context.Context, string, uint) error
	Debit(context.Context, string, uint) error
}

type Transaction interface {
	FindUnprocessed(ctx context.Context, limit int64) ([]transaction.DTO, error)
	GetOrderTransactions(context.Context, order.DTO) ([]transaction.DTO, error)
	RewardAccounts(ctx context.Context, limit int64) error
	Update(ctx context.Context, tx *transaction.DTO) error
}

type Repository struct {
	Account
	Order
	Transaction
}

func NewRepositoryMongoDB(db *mongo.Database) *Repository {
	return &Repository{
		Account:     mongodb.NewAccountStorage(db),
		Order:       mongodb.NewOrderStorage(db),
		Transaction: mongodb.NewTransactionStorage(db),
	}
}
