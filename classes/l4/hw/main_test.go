package main

import (
	"strings"

	"testing"

	"github.com/stretchr/testify/require"

	"hw/biglog" // You need to implement a biglog package.
)

// In this challenge, you need to implement a log system. It basically has 2 API, one for log and another for search.
func TestChallenge1(t *testing.T) {
	// 虽然这里说的是服务器，但是你不需要写任何网络相关的东西。这里的服务器是概念上的。
	server := biglog.StartServer() // starts a server and returns it.

	// The only requirement is that it accepts an io.Reader. You need to decide the storage mechanism.
	// 唯一的要求是 Log 函数接受一个 io.Reader。你自己需要决定储存机制。
	// Question:
	//		Should it accept io.ReadCloser instead?
	err := server.Log(strings.NewReader(`{"time": "2019-01-06", "number": 123}`))
	require.NoError(t, err) // Notice this log is the same.
	err = server.Log(strings.NewReader(`{"time": "2019-01-06", "number": 123}`))
	require.NoError(t, err) // Should you keep duplication?
	err = server.Log(strings.NewReader(`{"time": "2019-01-03", "number": 126}`))
	require.NoError(t, err)
	err = server.Log(strings.NewReader(`{"time": "2019-01-02", "number": 127}`))
	require.NoError(t, err)
	err = server.Log(strings.NewReader(`{"time": "2019-01-01", "number": 128, "price": 30}`))
	require.NoError(t, err)

	// // You should define a search API. It could have any signature.
	// // Find the record between 01-01 and 01-03 and return in logging order.
	// // Question:
	// //		Should it return errors? If so, what kind?
	var query biglog.QueryFunc
	logs, _ := server.All(
		query.Or(
			query.Key("time").Between("2019-01-02", "2019-01-03"),
			query.And(
				query.Key("number").Equal(128),
				query.Key("price").Exist(),
			),
		).End(),
	)
	s, _ := logs.ToStringArray()
	require.Equal(t, s, []string{
		`{"time":"2019-01-03","number":126}`, // Notice that JSON format is compact
		`{"time":"2019-01-02","number":127}`,
		`{"time": "2019-01-01","number":128,"price":30}`,
	})
}
