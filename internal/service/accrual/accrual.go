package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/pkg/httppc"
	"io"
	"net/http"
)

type Service struct {
	config config.AccrualServiceConfig
	client *httppc.Client
}

func (a *Service) GetAccrual(ctx context.Context, tx *transaction.Transaction) (uint, error) {
	urlString := fmt.Sprintf("%s/info?user=%s&good=%s&timestamp=%d",
		a.config.Endpoint, tx.UserID, tx.GoodID, tx.Timestamp)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		return 0, err
	}

	response, err := a.client.MakeRequest(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:

		bContent, err := io.ReadAll(response.Body)
		if err != nil {
			return 0, err
		}

		if err := json.Unmarshal(bContent, &tx.Reward); err != nil {
			return 0, err
		}
		return tx.Reward, nil

	case http.StatusTooManyRequests:
		return 0, apperr.AccrualServiceTooManyRequests
	case http.StatusNotFound:
		return 0, apperr.AccrualNotFound
	default:
		return 0, apperr.AccrualServiceInternalServerError
	}
}

func (a *Service) AdjustRPS(RPS int) {
	a.config.RPS = RPS
	a.client = httppc.New(a.config.MaxPoolSize, RPS)
}

func (a *Service) GetRPS() int {
	return a.config.RPS
}

func (a *Service) AdjustMaxPoolSize(MaxPoolSize int) {
	a.config.MaxPoolSize = MaxPoolSize
	a.client = httppc.New(a.config.MaxPoolSize, MaxPoolSize)
}

func New(config config.AccrualServiceConfig) *Service {
	return &Service{
		config: config,
		client: httppc.New(config.MaxPoolSize, config.RPS),
	}
}
