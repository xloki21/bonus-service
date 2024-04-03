package main

import (
	"context"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/application"
	"github.com/xloki21/bonus-service/pkg/log"
)

func main() {

	cfg, err := config.InitConfigFromViper()

	if err != nil {
		panic(err)
	}

	logger := log.GetDefaultLogger(cfg.LoggerConfig)
	app, err := application.New(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	if err := app.Run(context.Background()); err != nil {
		logger.Fatal(err)
	}
}
