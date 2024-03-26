package main

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/application"
	"log"
)

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/bonus-service/")
	viper.AddConfigPath("$HOME/.bonus-service")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("MONGO")
	if err := viper.BindEnv("USER"); err != nil {
		return err
	}

	if err := viper.BindEnv("PASSWORD"); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			viper.SetDefault("server.address", "localhost:8080")
			viper.SetDefault("store.uri", "mongodb://mongo-1:27017,mongo-2:27017,mongo-3:27017")
			viper.SetDefault("store.authdb", "admin")
			viper.SetDefault("store.dbname", "appdb")

			viper.SetDefault("transaction-service.polling_interval", "1000000000")
			viper.SetDefault("transaction-service.max_transactions_per_request", "10")
		} else {
			return err
		}
	}
	return nil
}

func main() {
	if err := initConfig(); err != nil {
		log.Fatalln("Failed to load config")
	}

	var cfg config.AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalln("Failed to unmarshal config")
	}
	cfg.DB.User = viper.GetString("user")
	cfg.DB.Password = viper.GetString("password")

	if cfg.DB.Password == "" || cfg.DB.User == "" {
		log.Fatalln("Failed to load credentials from config: set both user and password with environment variables")
	}

	app := application.New(&cfg)
	err := app.Run(context.Background())
	if err != nil {
		return
	}
}
