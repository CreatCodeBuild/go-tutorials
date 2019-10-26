package compute

// Operator is a any func that computes 2 integers
type Operator func(x, y int) int

func Add(x, y int) int {
	return x + y
}
func Sub(a, b int) int {
	return a - b
}
func Mul(x, y int) int {
	return x * y
}
func Div(x, y int) int {
	return x / y
}

func Compute(operants []int, operators []Operator) int {
	result := operants[0]
	for i := 0; i < len(operators); i++ {
		result = operators[i](result, operants[i+1])
	}
	return result
}
