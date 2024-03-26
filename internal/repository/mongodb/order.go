package mongodb

import (
	"context"
	"fmt"
	"github.com/xloki21/bonus-service/internal/apperr"
	t "github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type OrderMongoDB struct {
	db *mongo.Database
}

func NewOrderMongoDB(db *mongo.Database) *OrderMongoDB {
	return &OrderMongoDB{db: db}
}

// Register order to repository and create transactions
func (o *OrderMongoDB) Register(ctx context.Context, order *t.Order) error {
	var orders = o.db.Collection(ordersCollection)
	var transactions = o.db.Collection(transactionsCollection)
	filter := bson.D{
		{Key: "user_id", Value: order.UserID},
		{Key: "timestamp", Value: order.Timestamp},
	}
	var result t.Order

	if err := orders.FindOne(ctx, filter).Decode(&result); err == nil {
		return apperr.OrderAlreadyRegistered
	}

	// insert docs as single transaction
	_, err := o.Run(ctx, func(ctx context.Context) (interface{}, error) {
		// insert order into orders collection
		if _, err := orders.InsertOne(ctx, order); err != nil {
			return nil, fmt.Errorf("order registration failed: %w", err)
		}

		docs := make([]interface{}, 0, len(order.Goods))

		registeredAt := time.Now().Unix()
		for _, goodID := range order.Goods {
			tx := transaction.Transaction{
				UserID:       order.UserID,
				Status:       transaction.UNPROCESSED,
				GoodID:       goodID,
				Timestamp:    order.Timestamp,
				RegisteredAt: registeredAt,
			}
			docs = append(docs, tx)
		}
		// insert transactions into transactions collection
		if _, err := transactions.InsertMany(ctx, docs); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderMongoDB) Run(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	session, err := o.db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return f(sessCtx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
