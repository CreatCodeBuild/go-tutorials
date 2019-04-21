package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

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
			// There is an error in the producer. Shut down the mapping.
			// In some cases, you might want to keep going. We will discuss this scenario in later lectures.
			errs <- err
			return NewErrorChannel(errs, count+1)
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
		err := ConcurrentMapBetterErrorHandling(errorProducer(-1, 5), &results2, func(x interface{}) (interface{}, error) {
			return x, nil
		})
		switch e := err.(type) {
		case *ErrorChannel:
			require.Equal(t, 1, e.Count)
			require.EqualError(t, <-e.Errs, "producer error")
		case nil:
			t.FailNow()
		default:
			t.FailNow()
		}
	})

	t.Run("A mapper error followed by a producer error", func(t2 *testing.T) {
		results2 := outputConsumer2{}
		err := ConcurrentMapBetterErrorHandling(errorProducer(0, 5), &results2, func(x interface{}) (interface{}, error) {
			return nil, errors.New("mapper error")
		})
		switch e := err.(type) {
		case *ErrorChannel:
			require.Equal(t, 2, e.Count)
			// The order here is undefined. It's just because we only have 2 data here,
			// the second iteration of the Map happens to execute before the first goroutine starts.
			require.EqualError(t, <-e.Errs, "producer error")
			require.EqualError(t, <-e.Errs, "mapper error")
		case nil:
			t.FailNow()
		default:
			t.FailNow()
		}
	})

	t.Run("all 3 kinds of errors", func(t2 *testing.T) {
		err := ConcurrentMapBetterErrorHandling(errorProducer(1, 5), errorConsumer(), func(x interface{}) (interface{}, error) {
			i := x.(int)
			if i > 0 {
				return nil, errors.New("mapper error")
			}
			return nil, nil
		})
		switch e := err.(type) {
		case *ErrorChannel:
			require.Equal(t, 3, e.Count)
			fmt.Println(<-e.Errs)
			fmt.Println(<-e.Errs)
			fmt.Println(<-e.Errs)
		case nil:
			t.FailNow()
		default:
			t.FailNow()
		}
	})

	t.Run("many errors", func(t2 *testing.T) {
		err := ConcurrentMapBetterErrorHandling(errorProducer(998, 1000), errorConsumer(), func(x interface{}) (interface{}, error) {
			timeToSleep := time.Duration(rand.Int() % 3000)
			time.Sleep(timeToSleep * time.Millisecond)
			if rand.Int()%2 == 0 {
				return nil, fmt.Errorf("mapper error %v", x)
			}
			return nil, nil
		})
		switch e := err.(type) {
		case *ErrorChannel:
			require.Equal(t, 1000, e.Count)
			for i := 0; i < e.Count; i++ {
				fmt.Println(<-e.Errs)
			}
		case nil:
			t.FailNow()
		default:
			t.FailNow()
		}
	})
}

type SendFunc func(interface{}) error

func (f SendFunc) Send(i interface{}) error {
	return f(i)
}

func errorConsumer() SendFunc {
	return func(interface{}) error {
		return errors.New("consumer error")
	}
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

func errorProducer(after, n int) NextFunc {
	i := 0
	return func() (interface{}, error) {
		if i < n {
			defer func() { i++ }()
			if after < i {
				return i, errors.New("producer error")
			}
			return i, nil
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
