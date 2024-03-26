package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xloki21/bonus-service/config"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/pkg/httppc"
	"io"
	"net/http"
)

type AccrualServiceClient struct {
	config *config.AccrualServiceConfig
	client *httppc.Client
}

func (a *AccrualServiceClient) GetAccrual(ctx context.Context, tx *transaction.Transaction) (uint, error) {
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
	if response.StatusCode != http.StatusOK {
		return 0, errors.New("accrual service error")
	}
	bContent, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(bContent, &tx.Reward); err != nil {
		return 0, err
	}
	return tx.Reward, nil
}

func New(config *config.AccrualServiceConfig) *AccrualServiceClient {
	return &AccrualServiceClient{
		config: config,
		client: httppc.New(config.MaxPoolSize, config.RPS),
	}
}
