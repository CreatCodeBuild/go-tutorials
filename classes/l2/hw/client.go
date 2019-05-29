package hw

import (
	"net/http"
	"time"
)

// Client  A Http Client
type Client struct {
	BaseURL string
	client  *http.Client
}

// Request  request info of Client
type Request struct {
	Method, Path string
	client       *Client
	times        int
	timeout      time.Duration
}

// NewClient   Create a new Client
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, client: &http.Client{}}
}

// Request   Set http method & request path return a *Request
func (c *Client) Request(method string, path string) *Request {
	return &Request{Method: method, Path: path}
}

// WithArgs   Set params of the request
func (r *Request) WithArgs(args map[string]string) *Request {
	return r
}

// Do   Execute http request
func (c *Client) Do(request *Request) (*http.Response, error) {

	return nil, nil
}
