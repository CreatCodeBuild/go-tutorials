package hw

import (
	"net/http"
	"strings"
)

// Client  A Http Client
type Client struct {
	BaseURL string
	client  *http.Client
}

// Request  request info of Client
type Request struct {
	Method, Path string
	Args  map[string]string
}

// NewClient   Create a new Client
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, client: &http.Client{}}
}

// Request   Set http method & request path return a *Request
func (c *Client) Request(method string, path string) *Request {
	return &Request{Method: method, Path: path}
}

// WithArgs   Set args of the request
func (r *Request) WithArgs(args map[string]string) *Request {
	r.Args = args
	return r
}

// Do   Execute http request
func (c *Client) Do(request *Request) (*http.Response, error) {
	path := request.Path
	for k, v := range request.Args{
		placeholderStr := "{"+ k +"}"
		path = strings.Replace(path, placeholderStr,v,1)
	}
	path = c.BaseURL + path
	httpRequest,err := http.NewRequest(request.Method, path, nil)
	if err != nil {
		return nil, err
	}
	res,err := c.client.Do(httpRequest)
	return res, err
}
