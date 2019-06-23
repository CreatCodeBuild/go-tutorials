package main

import "fmt"

type inter interface {
	run()
}

var _ inter = x{}
type x struct {
}

func (x x) run() {
	fmt.Println("x")
}

type y struct {
	inter
}

func (y y) run() {
	if y.inter != nil {
		y.inter.run()
	}
}

func f(i inter) {
	i.run()
}


func main() {
	f(x{})
	// f(y{})
	y := y{}
	// y.inter = x{} // super
	fmt.Println(y.inter);
	f(y)
}