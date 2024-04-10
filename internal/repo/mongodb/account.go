package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/xloki21/bonus-service/internal/apperr"
	t "github.com/xloki21/bonus-service/internal/entity/account"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountMongoDB struct {
	db *mongo.Database
}

// Create new account.
func (a *AccountMongoDB) Create(ctx context.Context, account t.Account) error {
	accounts := a.db.Collection(accountsCollection)
	filter := bson.D{bson.E{Key: "user_id", Value: account.ID}}
	result := new(t.Account)

	if err := accounts.FindOne(ctx, filter).Decode(result); err == nil {
		return apperr.AccountAlreadyExists
	}

	if _, err := accounts.InsertOne(ctx, account); err != nil {
		return fmt.Errorf("can't create account: %w", err)
	}
	return nil
}

// Delete the account.
func (a *AccountMongoDB) Delete(ctx context.Context, account t.Account) error {
	accounts := a.db.Collection(accountsCollection)
	filter := bson.D{{Key: "user_id", Value: account.ID}}

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
func (a *AccountMongoDB) FindByID(ctx context.Context, id t.UserID) (*t.Account, error) {
	accounts := a.db.Collection(accountsCollection)

	filter := bson.D{{Key: "user_id", Value: id}}

	result := new(t.Account)

	if err := accounts.FindOne(ctx, filter).Decode(result); err != nil {
		return nil, apperr.AccountNotFound
	}
	return result, nil
}

// Credit credits account.
func (a *AccountMongoDB) Credit(ctx context.Context, id t.UserID, value uint) error {
	accounts := a.db.Collection(accountsCollection)
	filter := bson.D{{Key: "user_id", Value: id}}

	result := accounts.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"balance": value}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return apperr.AccountNotFound
	}

	return result.Err()
}

// Debit debits account.
func (a *AccountMongoDB) Debit(ctx context.Context, id t.UserID, value uint) error {

	accounts := a.db.Collection(accountsCollection)

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

func NewAccountMongoDB(db *mongo.Database) *AccountMongoDB {
	return &AccountMongoDB{db: db}
}
