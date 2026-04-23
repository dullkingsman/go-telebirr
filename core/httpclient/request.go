package httpclient

import (
	"bytes"
	"net/http"
)

// Request helps in constructing HTTP requests.
type Request struct {
	Method  string
	Url     string
	Body    []byte
	Headers map[string]string
}

// NewRequestBuilder initializes and returns a new Request instance.
func NewRequestBuilder() *Request {
	return &Request{
		Headers: make(map[string]string),
	}
}

// SetMethod sets the HTTP Method for the request.
func (rb *Request) SetMethod(method string) *Request {
	rb.Method = method
	return rb
}

// SetURL sets the URL for the request.
func (rb *Request) SetURL(url string) *Request {
	rb.Url = url
	return rb
}

// AddHeader adds a header key-value pair to the request.
func (rb *Request) AddHeader(key, value string) *Request {
	rb.Headers[key] = value
	return rb
}

// SetBody sets the Body content for the request.
func (rb *Request) SetBody(body []byte) *Request {
	rb.Body = body
	return rb
}

// Build constructs and returns an *http.Request object.
func (rb *Request) Build() (*http.Request, error) {
	req, err := http.NewRequest(rb.Method, rb.Url, bytes.NewBuffer(rb.Body))

	if err != nil {
		return nil, err
	}

	for key, value := range rb.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
