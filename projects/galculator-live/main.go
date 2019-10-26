package main

import (
	// std
	"fmt"
	"bufio"
	"os"

	// my
	"./compute"
)

// func parse(input string) (operants, operators)

func userInput() string {
	fmt.Println("please type your input")
	buf := bufio.NewReader(os.Stdin)
	line, _ := buf.ReadBytes('\n') // ascii 128 byte 256
	return string(line)
}

func main() {
	// terminal -> string
	fmt.Println(userInput())
	// string -> operants, operators
	// input := "1 + 3 - 8 * 10 / 2"
	operants := []int{1, 3, 8, 10, 2}
	operators := []compute.Operator{compute.Add, compute.Sub, compute.Mul, compute.Div}
	fmt.Println(compute.Compute(operants, operators))
}
