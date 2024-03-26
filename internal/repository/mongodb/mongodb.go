package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ordersCollection       = "orders"
	transactionsCollection = "transactions"
	accountsCollection     = "accounts"
)

type Config struct {
	URI      string `yaml:"URI"`
	AuthDB   string `yaml:"authdb"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

var TestDBConfig = Config{
	URI:      "mongodb://mongo-1:27017,mongo-2:27017,mongo-3:27017",
	AuthDB:   "admin",
	User:     "user",
	Password: "pass",
	DBName:   "appdb_test",
}

func NewMongoDB(ctx context.Context, cfg Config) (*mongo.Database, func(ctx context.Context) error, error) {
	clientOpts := options.Client().ApplyURI(cfg.URI)

	clientOpts.SetAuth(options.Credential{
		AuthSource: cfg.AuthDB,
		Username:   cfg.User,
		Password:   cfg.Password,
	})
	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to mongodb: %s", err.Error())
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %s", err.Error())
	}

	return client.Database(cfg.DBName), func(ctx context.Context) error {
		collNames, err := client.Database(cfg.DBName).ListCollectionNames(ctx, bson.D{})
		if err != nil {
			return err
		}
		for _, collName := range collNames {
			if err := client.Database(cfg.DBName).Collection(collName).Drop(ctx); err != nil {
				return err
			}
		}
		return nil
	}, nil
}
