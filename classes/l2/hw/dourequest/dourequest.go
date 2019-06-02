package dourequest

import (
	"bytes"
	"io/ioutil"
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
	body       []byte
	timeout    time.Duration
	retryTimes int
	client     *http.Client
}

func NewRequest(path string) *request {
	r := &request{
		client: &http.Client{},
		path:   path,
		args:   make(map[string]string),
		query:  url.Values{},
	}
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

func (r *request) SetArgs(args map[string]string) *request {
	r.args = args
	return r
}

func (r *request) Arg(key string, value string) *request {
	r.args[key] = value
	return r
}

func (r *request) SetQuery(query url.Values) *request {
	r.query = query
	return r
}

func (r *request) Query(key string, value string) *request {
	r.query.Set(key, value)
	return r
}

func (r *request) Body(body []byte) *request {
	r.body = body
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
	path := getFillArgsPath(r.path, r.args)

	// Timeout
	r.client.Timeout = r.timeout * time.Millisecond

	// requestUrl
	requestUrl := BaseURL + path

	// Request
	req, err := http.NewRequest(r.method, requestUrl, bytes.NewReader(r.body))
	if err != nil {
		return nil, err
	}

	// Query
	// FIXME 这么写的话,当 path 中附带了 query 参数时，就会出错了
	rawQuery := getQueryURL(r.query.Encode(), r.args)
	req.URL.RawQuery = rawQuery

	// Do
	return r.doAndRetry(req)
}

func getQueryURL(queryEncode string, args map[string]string) string {
	if args == nil || len(args) < 1 {
		return queryEncode
	}
	queryEncode = strings.Replace(queryEncode, "%7B", "{", 1)
	queryEncode = strings.Replace(queryEncode, "%7D", "}", 1)
	rawQuery := getFillArgsPath(queryEncode, args)
	return rawQuery
}

// 获取填充 args 后的 path
func getFillArgsPath(path string, args map[string]string) string {
	if len(args) == 0 {
		return path
	}
	// TODO 有待优化
	for k, v := range args {
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
		//fmt.Println("Timeout RetryTimes:", i+1,
		//	" Reset Timeout:", r.client.Timeout.Nanoseconds()/time.Millisecond.Nanoseconds(), "ms")

		// ContentLength Retry 没有重新创建 request,而io.Reader已经读到了末尾,导致request认为 ContentLength 为0,需要重置 body
		// https://stackoverflow.com/questions/31337891/net-http-http-contentlength-222-with-body-length-0
		req.Body = ioutil.NopCloser(bytes.NewReader(r.body))
		res, err = r.client.Do(req)
		netErr, ok = err.(timeoutError)
	}
	return res, err
}
