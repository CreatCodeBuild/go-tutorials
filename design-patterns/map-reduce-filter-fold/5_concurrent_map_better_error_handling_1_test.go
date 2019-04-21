package main

import (
	"errors"
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

func NewErrorChannel(errs chan error, count int) error {
	errCh := &ErrorChannel{
		Errs: make(chan error, count),
	}
	for i := 0; i < count; i++ {
		e := <-errs
		if e != nil {
			errCh.Errs <- e
			errCh.Count++
		}
	}
	if errCh.Count == 0 {
		return nil
	}
	return errCh
}

func ConcurrentMapBetterErrorHandling(p genericProducer, c genericConsumer, mapper genericMapper) error {
	count := 0
	errs := make(chan error, 1)
	for {
		next, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break // There is no more elements in the producer.
			}
			errs <- err
			return NewErrorChannel(errs, count+1) // There is an error in the producer. Shut down the mapping.
		}
		count++
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
	return NewErrorChannel(errs, count)
}

func TestConcurrentMapBetterErrorHandling(t *testing.T) {
	t.Run("all correct", func(t2 *testing.T) {
		results2 := outputConsumer2{}
		err := ConcurrentMapBetterErrorHandling(NewIntProducer(1, 2, 3, 4, 5), &results2, func(x interface{}) (interface{}, error) {
			return x, nil
		})
		require.NoError(t, err)
	})

	t.Run("all producer errors", func(t2 *testing.T) {
		results2 := outputConsumer2{}
		err := ConcurrentMapBetterErrorHandling(errorProducer(5), &results2, func(x interface{}) (interface{}, error) {
			return x, nil
		})
		require.NoError(t, err)
	})
}

// type genericProducer interface {
// 	Next() (interface{}, error)
// }

// NextFunc is a function type which implements a one method interface that has the same signature.
type NextFunc func() (interface{}, error)

// Next calls the function itself.
func (f NextFunc) Next() (interface{}, error) {
	return f()
}

func errorProducer(n int) NextFunc {
	i := 0
	return func() (interface{}, error) {
		if i < n {
			defer func() { i++ }()
			return nil, errors.New("some random error")
		}
		return nil, io.EOF
	}
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

// MultiProducer merges multiple producers into one. You can find the same design in io.MultiReader and io.MultiWriter, though their implementations are different.
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
