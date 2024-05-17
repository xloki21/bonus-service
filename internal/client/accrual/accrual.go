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

type Client struct {
	config     config.AccrualServiceConfig
	httpClient *httppc.Client
}

func (c *Client) SetHTTPClient(client httppc.HTTPDoer) {
	c.httpClient.SetClient(client)
}

func (c *Client) GetAccrual(ctx context.Context, tx *transaction.DTO) (uint, error) {
	urlString := fmt.Sprintf("%s/info?user=%s&good=%s&timestamp=%d",
		c.config.Endpoint, tx.UserID, tx.GoodID, tx.Timestamp)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		return 0, err
	}

	response, err := c.httpClient.MakeRequest(request)
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

func (c *Client) AdjustRPS(RPS int) {
	c.config.RPS = RPS
	c.httpClient = httppc.New(c.config.MaxPoolSize, RPS)
}

func (c *Client) GetRPS() int {
	return c.config.RPS
}

func (c *Client) AdjustMaxPoolSize(MaxPoolSize int) {
	c.config.MaxPoolSize = MaxPoolSize
	c.httpClient = httppc.New(c.config.MaxPoolSize, MaxPoolSize)
}

func NewClient(config config.AccrualServiceConfig) *Client {
	return &Client{
		config:     config,
		httpClient: httppc.New(config.MaxPoolSize, config.RPS),
	}
}
