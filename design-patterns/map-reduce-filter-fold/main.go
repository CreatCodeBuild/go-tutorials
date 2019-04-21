package main

import "fmt"

func main() {
	c := make(chan bool)
	close(c)
	select {
	case t := <-c:
		fmt.Println(t)
	}
}
