package homework

type expression interface {
	value() int
}

type operation struct {
	operator rune
	left     expression
	right    expression
}

func (op operation) value() int {
	switch op.operator {
	case '+':
		return op.left.value() + op.right.value()
	case '-':
		return op.left.value() + op.right.value()
	case '*':
		return op.left.value() * op.right.value()
	case '/':
		return op.left.value() / op.right.value()
	default:
		panic(string(op.operator) + " operator is invalid")
	}
}

type number int

func (n number) value() int {
	return int(n)
}

func LeftFirstTree(input string) expression {

}

func RightFirstTree(input string) expression {

}
