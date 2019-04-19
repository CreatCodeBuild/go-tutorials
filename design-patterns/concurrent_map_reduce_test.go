package main

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Map(data []int, mapper func(int) int) []int {
	results := make([]int, len(data))
	for i, ele := range data {
		results[i] = mapper(ele)
	}
	return results
}

// []int not a sub-type of []interface
// func(int) string is not a sub-type of func(interface{}) interface{}
// Reflection
func ReflectiveMap(data interface{}, mapper interface{}) interface{} {
	return nil
}

type producer interface {
	Next() interface{}
}

type consumer interface {
	Send(interface{}) error
}

// Object Oriented Design Tricks
func OOMap(p producer, c consumer, mapper func(interface{}) (interface{}, error)) error {
	for {
		next := p.Next()
		if next == nil {
			return nil
		}
		datum, err := mapper(next)
		if err != nil {
			return err
		}
		err = c.Send(datum)
		if err != nil {
			return err
		}
	}
}

type IntProducer struct {
	index int
	data  []int
}

func (ip *IntProducer) Next() interface{} {
	if ip.index < len(ip.data) {
		defer func() { ip.index++ }()
		return ip.data[ip.index]
	}
	return nil
}

type StringConsumer []string

func (sc *StringConsumer) Send(ele interface{}) error {
	if s, ok := ele.(string); ok {
		*sc = append(*sc, s)
		return nil
	}
	return errors.New("not a string")
}

func TestMapReduce(t *testing.T) {
	// Map
	results := Map([]int{1, 2, 3}, func(x int) int {
		return x + 1
	})
	require.Equal(t, []int{2, 3, 4}, results)

	// // ReflectiveMap Map
	// results2 := ReflectiveMap([]int{1, 2, 3}, func(x int) string {
	// 	return strconv.FormatInt(int64(x), 2)
	// })
	// require.Equal(t, []int{2, 3, 4}, results2)

	// Object Oriented Map
	// We hoist/lift the type from a slice to a producer without changing the structure.
	results2 := StringConsumer{}
	err := OOMap(&IntProducer{data: []int{1, 2, 3}}, &results2, func(x interface{}) (interface{}, error) {
		if i, ok := x.(int); ok {
			return strconv.FormatInt(int64(i), 2), nil
		}
		return nil, errors.New("lambda: not an int")
	})
	require.NoError(t, err)
	require.Equal(t, []string{"1", "10", "11"}, []string(results2))

	// Reduce / Unfold

	// Filter

	// Fold
}

func TestConcurrentMapReduce(t *testing.T) {

	// Map

	// Reduce / Unfold

	// Filter

	// Fold
}
