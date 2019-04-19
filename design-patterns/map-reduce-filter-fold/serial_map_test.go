package main

// import (
// 	"errors"
// 	"strconv"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// // ConcurrentMap
// func ConcurrentMap(p producer, c consumer, mapper func(interface{}) (interface{}, error)) error {
// 	for {
// 		next := p.Next()
// 		if next == nil {
// 			return nil
// 		}
// 		var errCh chan error
// 		var done chan struct{}
// 		go func(errChan chan error) {
// 			datum, err := mapper(next)
// 			if err != nil {
// 				errCh <- err
// 			}
// 			err = c.Send(datum)
// 			if err != nil {
// 				errCh <- err
// 			}
// 		}(errCh)

// 	}
// 	close(done)
// 	select {
// 	case err := <-errCh:
// 		return err
// 	case <-done:
// 		return nil
// 	}
// }

// func TestConcurrentMap(t *testing.T) {
// 	results2 := StringConsumer{}
// 	err := OOMap(&IntProducer{data: []int{1, 2, 3}}, &results2, func(x interface{}) (interface{}, error) {
// 		if i, ok := x.(int); ok {
// 			return strconv.FormatInt(int64(i), 2), nil
// 		}
// 		return nil, errors.New("lambda: not an int")
// 	})
// 	require.NoError(t, err)
// 	require.Equal(t, []string{"1", "10", "11"}, []string(results2))
// }
