package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response represents the outcome of an HTTP request, encapsulating
// status code, headers, body, and any related error.
type Response[T any] struct {
	Status  int
	Headers http.Header
	Body    T
}

func NewResponse[T any](resp *http.Response, err error) (*Response[T], error) {
	if err != nil {
		return nil, err
	}

	var tmp = &Response[T]{
		Status:  resp.StatusCode,
		Headers: resp.Header,
	}

	if err = json.NewDecoder(resp.Body).Decode(&tmp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return tmp, fmt.Errorf("HTTP error! status: %d", resp.StatusCode)
	}

	return tmp, nil
}
