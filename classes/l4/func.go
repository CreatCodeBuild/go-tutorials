package main

import "fmt"

type Int int

func (i Int) Apply(f func (int) int) int {
	return f(int(i))
}

func main() {
	i := Int(10)
	ret := i.Apply(func(x int)  int{
		return x * x
	})
	fmt.Println(ret)
}