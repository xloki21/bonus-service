package main

import (
	"context"
	"fmt"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/application"
	"github.com/xloki21/bonus-service/internal/pkg/log"
)

func main() {

	cfg, err := config.InitConfigFromViper()
	if err != nil {
		panic(err)
	}
	logger := log.GetDefaultLogger(cfg.LoggerConfig)
	app, err := application.New(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	if err := app.Run(context.Background()); err != nil {
		fmt.Println("````1231231231313123123123123")
		logger.Fatal(err)
	}
}
