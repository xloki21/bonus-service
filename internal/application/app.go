package application

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/config"
	controller "github.com/xloki21/bonus-service/internal/controller/http"
	v1 "github.com/xloki21/bonus-service/internal/controller/http/v1"
	"github.com/xloki21/bonus-service/internal/repository"
	"github.com/xloki21/bonus-service/internal/repository/mongodb"
	"github.com/xloki21/bonus-service/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	cfg      config.AppConfig
	repo     *repository.Repository
	services *service.Service
	teardown func(context.Context) error
	server   *controller.Server
}

func New(cfg config.AppConfig) (*Application, error) {
	db, teardown, err := mongodb.NewMongoDB(context.Background(), cfg.DB)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepositoryMongoDB(db)
	services := service.NewService(repo, cfg)
	return &Application{
		cfg:      cfg,
		repo:     repo,
		services: services,
		teardown: teardown,
		server:   &controller.Server{},
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	//a.logger.Info("Application started")
	defer func() {
		if err := a.teardown(ctx); err != nil {
			panic(err)
		}
	}()

	errCh := make(chan error, 1)

	go func() {
		err := a.services.Transaction.Polling(ctx)
		if err != nil {
			return
		}
	}()

	handler := v1.NewHandler(a.services)
	go func() {
		if err := a.server.Run(a.cfg.Server.Address, handler.ApiV1(a.cfg.Mode)); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		//a.logger.Info("Server shutdown signal received")
		return err
	case <-quit:
		//a.logger.Info("Gracefully shutting down...")
		return nil
	}
}
