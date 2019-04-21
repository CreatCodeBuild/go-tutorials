package main

import (
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func ConcurrentMapPoorErrorHandling(p genericProducer, c genericConsumer, mapper genericMapper) error {
	count := 0
	// empty struct{} is a type in Go. You can also redefine the type as:
	// type DoneSignal struct{}
	done := make(chan struct{})
	for {
		next, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break // There is no more elements in the producer.
			}
			return err // There is an error in the producer. Shut down the mapping.
		}
		count++
		go func(next interface{}) {
			ele, err := mapper(next)
			if err != nil {
				panic(err)
			}
			err = c.Send(ele)
			if err != nil {
				panic(err)
			}
			done <- struct{}{}
		}(next)
	}
	for i := 0; i < count; i++ {
		<-done
	}
	return nil
}

func TestConcurrentMap(t *testing.T) {
	results2 := outputConsumer2{}
	err := ConcurrentMapPoorErrorHandling(&intProducer{data: []int{1, 2, 3, 4, 5, 6, 7}}, &results2, func(x interface{}) (interface{}, error) {
		if i, ok := x.(int); ok {
			return strconv.FormatInt(int64(i), 2), nil
		}
		return nil, errors.New("lambda: not an int")
	})
	require.NoError(t, err)
}
