package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TransactionMongoDB struct {
	db *mongo.Database
}

func (t *TransactionMongoDB) GetOrderTransactions(ctx context.Context, order *order.Order) ([]transaction.Transaction, error) {
	transactions := t.db.Collection(transactionsCollection)
	var opts = options.Find()

	ops := make([]transaction.Transaction, 0, len(order.Goods))
	cursor, err := transactions.Find(ctx,
		bson.D{
			{Key: "user_id", Value: order.UserID},
			{Key: "timestamp", Value: order.Timestamp},
		}, opts)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var tx transaction.Transaction
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		ops = append(ops, tx)
	}
	return ops, nil
}

func (t *TransactionMongoDB) run(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	session, err := t.db.Client().StartSession()
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

// Update transaction status and reward
func (t *TransactionMongoDB) Update(ctx context.Context, tx *transaction.Transaction) error {
	transactions := t.db.Collection(transactionsCollection)
	filter := bson.M{"_id": tx.ID.(primitive.ObjectID)}
	update := bson.M{"$set": bson.M{"status": tx.Status, "reward": tx.Reward, "processed_at": tx.ProcessedAt}}
	_, err := transactions.UpdateOne(ctx, filter, update)
	return err
}

// RewardAccounts used to update accounts balance
func (t *TransactionMongoDB) RewardAccounts(ctx context.Context, limit int64) error {
	transactions := t.db.Collection(transactionsCollection)
	var opts = options.Aggregate()
	opts.SetBatchSize(int32(limit))
	aggTransactions := make([]transaction.AggregatedTransaction, 0, limit)

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "status",
						Value: bson.D{
							{Key: "$eq", Value: transaction.PROCESSED}}}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$user_id"},
					{Key: "reward", Value: bson.D{{Key: "$sum", Value: "$reward"}}},
					{Key: "transactions", Value: bson.D{{Key: "$push", Value: "$_id"}}}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "user_id", Value: "$_id"},
					{Key: "transactions", Value: "$transactions"},
					{Key: "reward", Value: "$reward"}}}}}

	cursor, err := transactions.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var op transaction.AggregatedTransaction
		if err := cursor.Decode(&op); err != nil {
			return err
		}
		aggTransactions = append(aggTransactions, op)
	}

	// Update accounts balances
	accounts := t.db.Collection(accountsCollection)
	_, err = t.run(ctx, func(ctx context.Context) (interface{}, error) {
		for _, tx := range aggTransactions {
			filter := bson.D{{Key: "user_id", Value: tx.UserID}}

			result := accounts.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"balance": tx.Reward}})
			if result.Err() != nil {
				return nil, fmt.Errorf("error during account balance update: %w", result.Err())
			}

			opRes, err := transactions.UpdateMany(ctx,
				bson.M{"_id": bson.M{"$in": tx.Transactions}},
				bson.M{"$set": bson.D{
					{Key: "status", Value: transaction.COMPLETED},
					{Key: "processed_at", Value: time.Now().Unix()}}})
			if err != nil {
				return nil, fmt.Errorf("error during transactions completion: %w", err)
			}
			if opRes.ModifiedCount != int64(len(tx.Transactions)) {
				return nil, errors.New("error during transactions completion: partial updates")
			}
		}
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionMongoDB) FindUnprocessed(ctx context.Context, limit int64) ([]transaction.Transaction, error) {
	transactions := t.db.Collection(transactionsCollection)
	var opts = options.Find()
	opts.SetSort(bson.D{{Key: "registered_at", Value: 1}})
	opts.SetLimit(limit)

	ops := make([]transaction.Transaction, 0, limit)
	cursor, err := transactions.Find(ctx, bson.D{
		{Key: "status", Value: bson.D{{Key: "$eq", Value: transaction.UNPROCESSED}}},
	}, opts)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var tx transaction.Transaction
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		ops = append(ops, tx)
	}
	return ops, nil
}

func NewTransactionMongoDB(db *mongo.Database) *TransactionMongoDB {
	return &TransactionMongoDB{db: db}
}
