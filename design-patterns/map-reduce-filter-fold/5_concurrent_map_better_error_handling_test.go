package main

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type ErrorChannel struct {
	Errs  chan error
	Count int
}

func (e *ErrorChannel) Error() string {
	return fmt.Sprintf("There are %d errors in this channel.", e.Count)
}

func ConcurrentMapBetterErrorHandling(p genericProducer, c genericConsumer, mapper genericMapper) error {
	count := 0
	errs := make(chan error)
	for {
		count++
		next, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break // There is no more elements in the producer.
			}
			errs <- err
			return &ErrorChannel{Errs: retErrs, Count: count}// There is an error in the producer. Shut down the mapping.
		}
		go func(next interface{}) {
			ele, err := mapper(next)
			if err != nil {
				errs <- err
				return
			}
			err = c.Send(ele)
			if err != nil {
				errs <- err
				return
			}
			errs <- nil
		}(next)
	}
	return &ErrorChannel{Errs: errs, Count: count}
}

// type genericProducer interface {
// 	Next() (interface{}, error)
// }

type NextFunc func() (interface{}, error)

func (f NextFunc) Next() (interface{}, error) {
	return f()
}

func NewIntProducer(slice ...int) NextFunc {
	pipe := make(chan int)
	done := make(chan struct{})
	go func() {
		for _, i := range slice {
			pipe <- i
		}
		close(done)
	}()
	return func() (interface{}, error) {
		select {
		case i := <-pipe:
			return i, nil
		case <-done:
			return 0, io.EOF
		}
	}
}

func MultiProducer(ps ...producer) NextFunc {
	type result struct {
		value interface{}
		error error
	}
	pipe := make(chan result)
	done := make(chan struct{})
	go func() {
		for _, p := range ps {
			next, err := p.Next()
			pipe <- result{next, err}
		}
		close(pipe)
	}()
	return func() (interface{}, error) {
		select {
		case r := <-pipe:
			return r.value, r.error
		case <-done:
			return 0, io.EOF
		}
	}
}

func TestConcurrentMapBetterErrorHandling(t *testing.T) {
	results2 := outputConsumer2{}
	err := ConcurrentMapBetterErrorHandling(NewIntProducer(1, 2, 3), &results2, func(x interface{}) (interface{}, error) {
		return x, nil
	})
	require.NoError(t, err)
}
