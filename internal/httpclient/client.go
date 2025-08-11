package httpclient

import (
	"github.com/dullkingsman/go-telebirr/internal/utils"
	"github.com/dullkingsman/go-telebirr/internal/values"
	"net/http"
	"time"
)

type RetryBackoffMethod[T any] func(attempt int, client *HTTPClient[T], resp *http.Response, err error) time.Duration

func DefaultRetryBackoffMethod[T any](attempt int, client *HTTPClient[T], resp *http.Response, err error) time.Duration {
	return time.Second * time.Duration(attempt) // simple backoff
}

type HTTPClient[T any] struct {
	client             *http.Client
	maxRetries         int
	retryBackoffMethod RetryBackoffMethod[T]
}

func NewHTTPClient[T any](config ...ClientConfig[T]) *HTTPClient[T] {
	var client = values.GetDefaultHttpClient()

	var _config = utils.GetOptionalParam(config...)

	var tmp = &HTTPClient[T]{
		client: _config.GetClient(client),
	}

	if tmp.client.Timeout == 0 {
		tmp.client.Timeout = time.Second * 10 // default (hoist later to defaults)
	}

	tmp.maxRetries = 0
	tmp.retryBackoffMethod = DefaultRetryBackoffMethod

	if _config != nil {
		tmp.Configure(_config)
	}

	return tmp
}

func (c *HTTPClient[T]) GetClient() *http.Client {
	return c.client
}

func (c *HTTPClient[T]) SetMaxRetries(maxRetries int) *HTTPClient[T] {
	if c == nil {
		return nil
	}

	c.maxRetries = maxRetries

	return c
}

func (c *HTTPClient[T]) GetMaxRetries() int {
	if c == nil {
		return 0
	}

	return c.maxRetries
}

func (c *HTTPClient[T]) SetRetryBackoffMethod(retryBackoffMethod RetryBackoffMethod[T]) *HTTPClient[T] {
	if c == nil {
		return nil
	}

	c.retryBackoffMethod = retryBackoffMethod
	return c
}

func (c *HTTPClient[T]) Configure(config *ClientConfig[T]) *HTTPClient[T] {
	if c == nil || config == nil {
		return nil
	}

	if config.MaxRetries != nil {
		c.SetMaxRetries(*config.MaxRetries)
	}

	if config.RetryBackoffMethod != nil {
		c.SetRetryBackoffMethod(config.RetryBackoffMethod)
	}

	return c
}

func (c *HTTPClient[T]) DoRequest(req *Request) (*Response[T], error) {
	var (
		resp *http.Response
		err  error
		_req *http.Request
	)

	if _req, err = req.Build(); err != nil {
		return nil, err
	}

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err = c.client.Do(_req)

		if err == nil {
			return NewResponse[T](resp, err)
		}

		time.Sleep(c.retryBackoffMethod(attempt, c, resp, err))
	}

	return nil, err
}
