package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountStorage struct {
	db *mongo.Database
}

func (a *AccountStorage) collection(name string) *mongo.Collection {
	return a.db.Collection(name)
}

// Create new account.
func (a *AccountStorage) Create(ctx context.Context, acc account.DTO) error {
	accounts := a.collection(accountsCollection)

	if _, err := accounts.InsertOne(ctx, acc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return apperr.AccountAlreadyExists
		}
		return fmt.Errorf("can't create account: %w", err)
	}
	return nil
}

// Delete the account.
func (a *AccountStorage) Delete(ctx context.Context, acc account.DTO) error {
	accounts := a.collection(accountsCollection)
	filter := bson.D{{Key: "user_id", Value: acc.ID}}

	opResult, err := accounts.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't delete account: %w", err)
	}
	if opResult.DeletedCount == 0 {
		return apperr.AccountNotFound
	}
	return nil
}

// FindByID finds account by user id.
func (a *AccountStorage) FindByID(ctx context.Context, id string) (*account.DTO, error) {
	accounts := a.collection(accountsCollection)

	filter := bson.D{{Key: "user_id", Value: id}}

	result := &account.DTO{}

	if err := accounts.FindOne(ctx, filter).Decode(result); err != nil {
		return nil, apperr.AccountNotFound
	}
	return result, nil
}

// Credit credits account.
func (a *AccountStorage) Credit(ctx context.Context, id string, value uint) error {
	accounts := a.collection(accountsCollection)
	filter := bson.D{{Key: "user_id", Value: id}}

	result := accounts.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"balance": value}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return apperr.AccountNotFound
	}
	if result.Err() != nil {
		return fmt.Errorf("can't credit account: %w", result.Err())
	}
	return nil
}

// Debit debits account.
func (a *AccountStorage) Debit(ctx context.Context, id string, value uint) error {

	accounts := a.collection(accountsCollection)

	filter := bson.D{
		{Key: "user_id", Value: id},
		{Key: "balance", Value: bson.M{"$gte": value}},
	}
	result := accounts.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"balance": -int(value)}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return apperr.InsufficientBalance
	}
	return nil

}

func NewAccountStorage(db *mongo.Database) *AccountStorage {
	return &AccountStorage{db: db}
}
