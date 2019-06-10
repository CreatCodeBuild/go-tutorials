package biglog

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Log struct {}

var content []string

func StartServer() *Log {
	return &Log{}
}

func (this *Log) Log(r io.Reader) error {
	p := make([]byte, 1024)
	byteLen, _ := r.Read(p)
	if byteLen <= 0 {
		return errors.New("接收参数为空")
	}

	content = append(content, string(p[:byteLen]))
	return nil
}

func (this *Log) Search(s interface{}) []string {
	fmt.Println("search")
	fmt.Println(content)
	return []string{}
}

func (this *Log) Close() error {
	os.Exit(0)

	return nil
}