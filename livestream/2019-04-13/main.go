package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	ret := &StringSeq{}
	Map(
		newIntSeqConcur([]int{16, 17, 18, 19, 20}),
		ret,
		func(i interface{}) interface{} {
			time.Sleep(time.Duration(rand.Int()%5) * time.Second)
			integer, ok := i.(int)
			if !ok {
				return errors.New("xx")
			}
			return strconv.FormatInt(int64(integer), 16)
		})
	fmt.Println(ret)
}

type IntSeq struct {
	i    int
	ints []int
}

func (iq *IntSeq) Next() interface{} {
	if iq.i < len(iq.ints) {
		defer func() {
			iq.i++
		}()
		return iq.ints[iq.i]
	}
	return nil
}

type IntSeqConcur struct {
	ch chan interface{}
}

func newIntSeqConcur(ints []int) IntSeqConcur {
	seq := IntSeqConcur{
		ch: make(chan interface{}),
	}
	go func(seq IntSeqConcur) {
		for _, interger := range ints {
			seq.ch <- interger
		}
		close(seq.ch)
	}(seq)
	return seq
}

func (iq IntSeqConcur) Next() interface{} {
	return <-iq.ch
}

type StringSeq struct {
	strings []string
}

func (ss *StringSeq) Send(obj interface{}) {
	s, ok := obj.(string)
	if !ok {
		panic("not a string")
	}
	ss.strings = append(ss.strings, s)
}

type sequence interface {
	Next() interface{}
}

type receiver interface {
	Send(interface{})
}

// Map s
func Map(
	sequence sequence,
	r receiver,
	lambda func(interface{}) interface{},
) {
	// waitG := sync.WaitGroup{}
	count := 0
	done := make(chan struct{})
	for next := sequence.Next(); next != nil; next = sequence.Next() {
		// waitG.Add(1)
		// go func(r receiver, next interface{}, wg *sync.WaitGroup) {
		// 	r.Send(lambda(next))
		// 	waitG.Done()
		// }(r, next, &waitG)
		count++
		go func(r receiver, next interface{}, done chan struct{}) {
			defer func() { done <- struct{}{} }()
			r.Send(lambda(next))
		}(r, next, done)
	}
	// waitG.Wait()
	for i := 0; i < count; i++ {
		<-done
	}
}

// type []int != type []interface{}
// {1, 2, 3, 4}

// Reduce

// Filter

// Unfold
