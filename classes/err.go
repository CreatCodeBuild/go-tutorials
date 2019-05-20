package main

import "errors"
import "time"
import "fmt"

type X struct {
	a int
	b int
	err chan error
}

func (x X) Clone() X {
	return X{
		err: x.err,
	}
}

func (x X ) f() (int) {
	x.err <- errors.New("123")
	return 0
}

func (x X) Err() chan error {
	return x.err
}

func f1(x X) {
	go f3(x.Clone())
	x.err <- errors.New("Something bad happended!")
}

func f3(x X) {
	time.Sleep(time.Second * time.Duration(5))
	select {
	case <-x.Err():
		fmt.Println("My parent died. Release some resource")
	default:
		fmt.Println("Good")
	}
}

func main() {
	x := X{err: make(chan error)}
	f1(x)
	fmt.Println("Unhandled error", <-x.Err())
}
