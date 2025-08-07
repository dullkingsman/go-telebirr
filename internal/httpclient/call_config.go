package httpclient

import "net/http"

type ClientConfig[T any] struct {
	Client             *http.Client
	MaxRetries         *int
	RetryBackoffMethod RetryBackoffMethod[T]
}

func (c *ClientConfig[T]) GetClient(fallback *http.Client) *http.Client {
	if c == nil {
		return fallback
	}

	if c.Client != nil {
		return c.Client
	}

	return fallback
}
