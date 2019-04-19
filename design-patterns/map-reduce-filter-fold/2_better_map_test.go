package main

import (
	"fmt"
	"io"
	"strconv"
)

// producer 是一个数据生产者。Next 会迭代并返回序列中的下一个元素。
// 返回 io.EOF 表示穷尽了序列。其他错误值表示 producer 本身遇到了错误。
type producer interface {
	Next() (string, error)
}

// consumer 是数据消费者。Send 会读入新的数据。
type consumer interface {
	Send(int64)
}

// 返回错误如果 string 不能表示 int。比如 "xxx" 不是一个正确的 int 表示形式。
type mapper func(string) (int64, error)

func BetterMap(p producer, c consumer, mapper mapper) error {
	for {
		next, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err // 生产者本身遇到错误，终止 Map。
			}
		}
		datum, err := mapper(next)
		if err != nil {
			return err // mapper 出了问题，终止 Map。
		}
		c.Send(datum)
	}
	return nil
}

type StringProducer struct {
	index int
	data  []string
}

func (ip *StringProducer) Next() (string, error) {
	if ip.index < len(ip.data) {
		defer func() { ip.index++ }()
		return ip.data[ip.index], nil
	}
	return "", io.EOF
}

type OutputConsumer struct{}

func (c OutputConsumer) Send(ele int64) {
	fmt.Println(ele)
}

func ExampleBetterMap() {
	BetterMap(&StringProducer{data: []string{"1", "10", "11"}}, OutputConsumer{}, func(str string) (int64, error) {
		// 这里的 lambda 将字符串以二进制形式转为整数
		return strconv.ParseInt(str, 2, 64)
	})
	// Output: 1
	// 2
	// 3
}
