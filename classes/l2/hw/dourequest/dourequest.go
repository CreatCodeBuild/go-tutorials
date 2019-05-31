package dourequest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var BaseURL = ""

type timeoutError interface{ Timeout() bool }

type request struct {
	//baseURL string
	method     string
	path       string
	args       map[string]string
	query      url.Values
	timeout    time.Duration
	retryTimes int
	client     *http.Client
}

func NewRequest(path string) *request {
	r := &request{client: &http.Client{}, path: path}
	r.Timeout(1000)
	r.RetryTimes(3)
	return r
}

func Get(path string) *request {
	return NewRequest(path).Method(http.MethodGet)
}

func (r *request) Method(method string) *request {
	r.method = method
	return r
}

func (r *request) WithArgs(args map[string]string) *request {
	r.args = args
	return r
}

func (r *request) Query(query url.Values) *request {
	r.query = query
	return r
}

// 单位 毫秒
func (r *request) Timeout(timeout time.Duration) *request {
	r.timeout = timeout
	return r
}

func (r *request) RetryTimes(times int) *request {
	r.retryTimes = times
	return r
}

func (r *request) Do() (*http.Response, error) {
	// WithArgs
	path := r.getFillArgsPath()

	// Timeout
	r.client.Timeout = r.timeout * time.Millisecond

	// url
	url := BaseURL + path

	// Request
	req, err := http.NewRequest(r.method, url, nil)
	if err != nil {
		return nil, err
	}

	// Query
	req.URL.RawQuery = r.query.Encode()

	// Do
	return r.doAndRetry(req)
}

// 获取填充 args 后的 path
func (r *request) getFillArgsPath() string {
	if len(r.args) == 0 {
		return r.path
	}
	path := r.path
	for k, v := range r.args {
		placeholderStr := "{" + k + "}"
		path = strings.Replace(path, placeholderStr, v, 1)
	}
	return path
}

// 超时重试的请求
func (r *request) doAndRetry(req *http.Request) (*http.Response, error) {
	res, err := r.client.Do(req)
	netErr, ok := err.(timeoutError)
	for i := 0; ok && netErr.Timeout() && i < r.retryTimes; i++ {
		r.client.Timeout = r.client.Timeout * 2
		fmt.Println("Timeout RetryTimes:", i+1,
			" Reset Timeout:", r.client.Timeout.Nanoseconds()/time.Millisecond.Nanoseconds(), "ms")
		res, err = r.client.Do(req)
		netErr, ok = err.(timeoutError)
	}
	return res, err
}
