package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/xloki21/bonus-service/internal/repo/mongodb"
	"os"
	"time"
)

type Server struct {
	Address string `json:"address"`
}

type AccrualServiceConfig struct {
	Endpoint    string `json:"endpoint" mapstructure:"endpoint"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
	RPS         int    `json:"rps" mapstructure:"rps"`
}

type TransactionServiceConfig struct {
	PollingInterval           time.Duration `mapstructure:"polling_interval"`
	MaxTransactionsPerRequest int           `mapstructure:"max_transactions_per_request"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

type AppConfig struct {
	Mode                     string                   `mapstructure:"mode"`
	Server                   Server                   `mapstructure:"server"`
	DB                       mongodb.Config           `mapstructure:"store"`
	AccrualService           AccrualServiceConfig     `mapstructure:"accrual-service"`
	TransactionServiceConfig TransactionServiceConfig `mapstructure:"transaction-service"`
	LoggerConfig             LoggerConfig             `mapstructure:"logger"`
}

func InitConfigFromViper() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/bonus-service/")
	viper.AddConfigPath("$HOME/.bonus-service")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("BS") // BONUS_SERVICE prefix
	if err := viper.BindEnv("USER"); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("PASSWORD"); err != nil {
		return nil, err
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
			viper.SetDefault("logger.level", "info")
			viper.SetDefault("logger.encoding", "json")

		} else {
			return nil, err
		}
	}

	cfg := &AppConfig{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if _, ok := os.LookupEnv("BS_USER"); !ok {
		return nil, errors.New("missing environment variable BS_USER")
	}

	if _, ok := os.LookupEnv("BS_PASSWORD"); !ok {
		return nil, errors.New("missing environment variable BS_PASSWORD")
	}

	cfg.DB.User = viper.GetString("user")
	cfg.DB.Password = viper.GetString("password")
	return cfg, nil
}
