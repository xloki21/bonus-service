package config

import (
	"github.com/xloki21/bonus-service/internal/repository/mongodb"
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

type AppConfig struct {
	Mode                     string                   `mapstructure:"mode"`
	Server                   Server                   `mapstructure:"server"`
	DB                       mongodb.Config           `mapstructure:"store"`
	AccrualService           AccrualServiceConfig     `mapstructure:"accrual-service"`
	TransactionServiceConfig TransactionServiceConfig `mapstructure:"transaction-service"`
}
