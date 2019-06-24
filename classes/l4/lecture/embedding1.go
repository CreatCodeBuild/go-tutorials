package main

import "fmt"

// 别名
type Int3 = int // typedef

// Definition 类型重定义
type Int2 int

type Int1 struct {
	int // embeding 嵌入
}

type x struct {}

func (x x) run() {
	fmt.Println("x")
}

type y struct {
	x
}

func (y y) run() {
	fmt.Println("y")
	y.x.run()
}

type z struct {
	x
	y
}

func (z z) run() {
	fmt.Println("z")
	z.y.run()
}

type inter interface {
	run()
}

func f(i inter) {
	i.run()
}

func main() {
	f(z{})
}