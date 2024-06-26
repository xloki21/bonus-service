package application

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/config"
	controller "github.com/xloki21/bonus-service/internal/controller/http"
	v1 "github.com/xloki21/bonus-service/internal/controller/http/v1"
	"github.com/xloki21/bonus-service/internal/repo"
	"github.com/xloki21/bonus-service/internal/repo/mongodb"
	"github.com/xloki21/bonus-service/internal/service"
	"github.com/xloki21/bonus-service/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	cfg      config.AppConfig
	repo     *repo.Repository
	services *service.Service
	teardown func(context.Context) error
	server   *controller.Server
}

func New(cfg config.AppConfig) (*Application, error) {
	db, teardown, err := mongodb.NewMongoDB(context.Background(), cfg.DB)
	if err != nil {
		return nil, err
	}

	repos := repo.NewRepositoryMongoDB(db)
	services := service.NewService(repos, cfg)
	return &Application{
		cfg:      cfg,
		repo:     repos,
		services: services,
		teardown: teardown,
		server:   &controller.Server{},
	}, nil
}

func (a *Application) Run(ctx context.Context) (err error) {
	var logger log.Logger
	logger, err = log.GetLogger()
	if err != nil {
		return err
	}
	logger.Info("Application started")
	defer func() {
		err = a.teardown(ctx)
	}()

	errCh := make(chan error, 1)

	go func() {
		err = a.services.Transaction.Polling(ctx)
		if err != nil {
			return
		}
	}()

	handler := v1.NewHandler(a.services)
	go func() {
		if err = a.server.Run(a.cfg.Server.Address, handler.ApiV1(a.cfg.Mode)); err != nil {
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
	case err = <-errCh:
		return err
	case <-quit:
		return nil
	}
}
