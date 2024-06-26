package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TransactionStorage struct {
	db *mongo.Database
}

func (t *TransactionStorage) collection(name string) *mongo.Collection {
	return t.db.Collection(name)
}

func (t *TransactionStorage) GetOrderTransactions(ctx context.Context, order order.DTO) ([]transaction.DTO, error) {
	transactions := t.collection(transactionsCollection)
	opts := options.Find()

	ops := make([]transaction.DTO, 0, len(order.Goods))
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
		var tx transaction.DTO
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		ops = append(ops, tx)
	}
	return ops, nil
}

func (t *TransactionStorage) run(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
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

// Update transaction status and reward.
func (t *TransactionStorage) Update(ctx context.Context, tx *transaction.DTO) error {
	transactions := t.collection(transactionsCollection)
	filter := bson.M{"_id": tx.ID.(primitive.ObjectID)}
	update := bson.M{"$set": bson.M{"status": tx.Status, "reward": tx.Reward, "completed_at": tx.CompletedAt}}
	_, err := transactions.UpdateOne(ctx, filter, update)
	return err
}

// RewardAccounts used to update accounts balance.
func (t *TransactionStorage) RewardAccounts(ctx context.Context, limit int64) error {
	transactions := t.collection(transactionsCollection)
	opts := options.Aggregate()
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
	accounts := t.collection(accountsCollection)
	_, err = t.run(ctx, func(ctx context.Context) (interface{}, error) {
		for _, tx := range aggTransactions {
			filter := bson.D{{Key: "user_id", Value: tx.UserID}}

			result := accounts.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"balance": tx.Reward}})
			if result.Err() != nil {
				if errors.Is(result.Err(), mongo.ErrNoDocuments) {
					return nil, fmt.Errorf("account #{%s} balance update error: %w", filter[0].Value, apperr.AccountNotFound)
				}
				return nil, fmt.Errorf("account #{%s} balance update error: %w", filter[0].Value, result.Err())
			}

			opRes, err := transactions.UpdateMany(ctx,
				bson.M{"_id": bson.M{"$in": tx.Transactions}},
				bson.M{"$set": bson.D{
					{Key: "status", Value: transaction.COMPLETED},
					{Key: "completed_at", Value: time.Now().Unix()}}})
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

// FindUnprocessed returns unprocessed transactions.
func (t *TransactionStorage) FindUnprocessed(ctx context.Context, limit int64) ([]transaction.DTO, error) {
	transactions := t.collection(transactionsCollection)
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "registered_at", Value: 1}})
	opts.SetLimit(limit)

	ops := make([]transaction.DTO, 0, limit)
	cursor, err := transactions.Find(ctx, bson.D{
		{Key: "status", Value: bson.D{{Key: "$eq", Value: transaction.UNPROCESSED}}},
	}, opts)

	if err != nil {
		return nil, fmt.Errorf("find unprocessed transactions error: %w", err)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var tx transaction.DTO
		if err := cursor.Decode(&tx); err != nil {
			return nil, fmt.Errorf("find unprocessed transactions error: %w", err)
		}
		ops = append(ops, tx)
	}
	return ops, nil
}

func NewTransactionStorage(db *mongo.Database) *TransactionStorage {
	return &TransactionStorage{db: db}
}
