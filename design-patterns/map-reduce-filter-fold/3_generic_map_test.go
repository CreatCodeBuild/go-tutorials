package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// producer 是一个数据生产者。Next 会迭代并返回序列中的下一个元素。
// 返回 io.EOF 表示穷尽了序列。其他错误值表示 producer 本身遇到了错误。
type genericProducer interface {
	Next() (interface{}, error)
}

// consumer 是数据消费者。Send 会读入新的数据。
type genericConsumer interface {
	Send(interface{}) error
}

type genericMapper func(interface{}) (interface{}, error)

func GenericMap(p genericProducer, c genericConsumer, mapper genericMapper) error {
	for {
		next, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err // 生产者本身遇到错误，终止 Map。 The producer has an error. Shut down Map.
			}
		}
		datum, err := mapper(next)
		if err != nil {
			return err // mapper 出了问题，终止 Map。The mapper has an error. Shut down Map.
		}
		err = c.Send(datum)
		if err != nil {
			return err // Shut down the Map immediately. This is a design decision you have to make. You can also pipe the error to a channel and collect all errors.
		}
	}
	return nil
}

type intProducer struct {
	index int
	data  []int
}

func (ip *intProducer) Next() (interface{}, error) {
	if ip.index < len(ip.data) {
		defer func() { ip.index++ }()
		return ip.data[ip.index], nil
	}
	return 0, io.EOF
}

type outputConsumer2 struct{}

func (c outputConsumer2) Send(ele interface{}) error {
	fmt.Println("outputConsumer2", ele)
	return nil
}

func ExampleGenericMap() {
	GenericMap(&intProducer{data: []int{10, 11, 12}}, outputConsumer2{}, func(ele interface{}) (interface{}, error) {
		// 这里的 lambda 将字符串以二进制形式转为整数
		if i, ok := ele.(int); ok {
			return strconv.FormatInt(int64(i), 16), nil
		}
		return nil, errors.New("mapper: not an int")
	})
	// Output: a
	// b
	// c
}
