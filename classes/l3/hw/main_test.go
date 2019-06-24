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
	err := server.Log(strings.NewReader(`{"time": "2019-01-06", "number": 123}`)); require.NoError(t, err)	// Notice this log is the same.
	err = server.Log(strings.NewReader(`{"time": "2019-01-06", "number": 123}`)); require.NoError(t, err)	// Should you keep duplication?
	err = server.Log(strings.NewReader(`{"time": "2019-01-03", "number": 126}`)); require.NoError(t, err)
	err = server.Log(strings.NewReader(`{"time": "2019-01-02", "number": 127}`)); require.NoError(t, err)
	err = server.Log(strings.NewReader(`{"time": "2019-01-01", "number": 128}`)); require.NoError(t, err)

	// // You should define a search API. It could have any signature.
	// // Find the record between 01-01 and 01-03 and return in logging order.
	// // Question:
	// //		Should it return errors? If so, what kind?
	logs, _ := server.All(func(record biglog.Record) bool { 
		return record["time"].(string) == "2019-01-03"
	})
	require.Equal(t, logs, []string{
		`{"time":"2019-01-03","number":126}`,  // Notice that JSON format is compact
		`{"time":"2019-01-02","number":127}`,
		`{"time":"2019-01-01","number":128}`,
	})

	// // Find the record between 01-01 and 01-03, select 'number' only and return in logging order
	// logs := server.Search(?)
	// require.Equal(t, logs, []string{
	// 	`{"number":126}`,
	// 	`{"number":127}`,
	// 	`{"number":128}`
	// })

	// // Question:
	// //		Because this is a log system, should you trust the 'time' which users input? Should the system have its own time?
	// //		因为这是一个日志系统，你应该相信用户输入的时间吗？系统需要自己记录时间吗？ 

	// err = server.Close()
	// require.NoError(t, err)

	// // After the server is closed, it will refuse any api.
	// err = server.Log("does not matter")
	// // What should the err be?

	// ? = server.Search(?) // So, do you think Search should return an error now?

	// Question:
	//		这里我们的测试有一个明显的问题，就是没有并发测试。比如说，如果同时有 Log 和 Close 的调用会如何？
	//		你能够补全一些并发测试，来找到已有实现的问题吗？
	//		这一题为选做加分项，也是下一节课的内容。
}
