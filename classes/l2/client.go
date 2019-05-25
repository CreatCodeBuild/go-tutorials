package main

import "net/http"
import "fmt"
import "strconv"
import "net/url"
import "io/ioutil"
import "encoding/json"
import "time"

const (
	API1 = "http://localhost:8081"
	// API2 = "http://localhost:8080"
)

func main() {

	// client := http.Client{Timeout: time.Second}
	// fmt.Println(getUser(client, 10, 18, true))

	client2 := Client{client: &http.Client{}}
	client2.baseURL = API1
	fmt.Println(getUser2(client2,
		map[string]interface{}{"id": 10},
		map[string]interface{}{"gender": true, "age": 18},
	))

	c3, err := client.Subclient("/user/{id}")
	c3.SetHeader("jwt", "xxx")
	c3.Request().SetArg("id", 123)

	c4 := c3.Subclient("/orders")
}

type Client struct {
	client  *http.Client
	baseURL string
}

func (c *Client) Request(method, path string) *Request {
	fmt.Println(c, c.client)
	return &Request{method: method, path: path, client: c}
}

type Request struct {
	method, path string
	client       *Client
	times        int
	timeout      time.Duration
}

func (r *Request) Retry(times int) *Request {
	r.times = times
	return r
}

func (r *Request) Timeout(timeout time.Duration) *Request {
	r.timeout = timeout
	return r
}

func (r *Request) Do() (*http.Response, error) {
	path := r.client.baseURL + r.path
	req, err := http.NewRequest(r.method, path, nil)
	if err != nil {
		return nil, err
	}

	r.client.client.Timeout = time.Second * r.timeout
	res, err := r.client.client.Do(req)
	type timeout interface{ Timeout() bool }
	netErr, ok := err.(timeout)
	for i := 0; ok && netErr.Timeout() && i < r.times; i++ {
		fmt.Println("!")
		res, err = r.client.client.Do(req)
		r.client.client.Timeout = r.client.client.Timeout * r.timeout(i+1)
		netErr, ok = err.(timeout)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

type User struct {
	Age    int
	Gender bool
	ID     int
}

func getUser(client http.Client, userID int, age int, gender bool) (User, error) {

	path := API1 + "/users/" + strconv.FormatInt(int64(userID), 10) + "?"
	v := url.Values{}
	v.Add("age", strconv.FormatInt(int64(age), 10))
	v.Add("gender", strconv.FormatBool(gender))
	path = path + v.Encode()

	//
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	res, err := client.Do(req)
	for i := 0; i < 3; i++ {
		if err != nil {
			client.Timeout = client.Timeout * 2
			res, err = client.Do(req)
		} else {
			break
		}
	}

	//
	bytes, _ := ioutil.ReadAll(res.Body)
	user := User{}
	json.Unmarshal(bytes, &user)
	return user, err
}

func getUser2(client Client, restArgs map[string]interface{}, urlArgs map[string]interface{}) (User, error) {
	_, err := client.Request(http.MethodGet, "/users").
		Retry(func(i int) { return 2 * i }).
		Timeout(2).Do()
	// bytes, _ := ioutil.ReadAll(res.Body)
	user := User{}
	// json.Unmarshal(bytes, &user)
	return user, err
}
