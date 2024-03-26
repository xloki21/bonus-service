package httppc

import (
	"net/http"
	"time"
)

type Client struct {
	client       *http.Client
	maxPoolSize  int
	sem          chan interface{}
	reqPerSecond int
	rateLimiter  <-chan time.Time
}

func New(maxPoolSize int, reqPerSec int) *Client {
	var sem chan interface{} = nil
	if maxPoolSize > 0 {
		sem = make(chan interface{}, maxPoolSize)
	}

	var emitter <-chan time.Time = nil
	if reqPerSec > 0 {
		emitter = time.NewTicker(time.Second / time.Duration(reqPerSec)).C
	}

	return &Client{
		client:       &http.Client{},
		maxPoolSize:  maxPoolSize,
		sem:          sem,
		reqPerSecond: reqPerSec,
		rateLimiter:  emitter,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.MakeRequest(req)
}

func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	if c.maxPoolSize > 0 {
		c.sem <- struct{}{}
		defer func() {
			<-c.sem
		}()
	}

	if c.reqPerSecond > 0 {
		<-c.rateLimiter
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil
}
