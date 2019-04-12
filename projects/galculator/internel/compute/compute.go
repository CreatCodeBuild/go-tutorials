package compute

import (
	"errors"
	"galculator/internel/lexer"
	"strconv"
)

type variableMap map[string]int64

// todo: refactor operator to operator expression
type operator func(left expression, right expression) (NumberExpression, error)

func add(left, right expression) (ne NumberExpression, err error) {
	lv, err2 := left.Value()
	if err2 != nil {
		err = err2
		return
	}
	rv, err2 := right.Value()
	if err2 != nil {
		err = err2
		return
	}
	return NumberExpression{lv + rv}, nil
}

func sub(left, right expression) (ne NumberExpression, err error) {
	lv, err2 := left.Value()
	if err2 != nil {
		err = err2
		return
	}
	rv, err2 := right.Value()
	if err2 != nil {
		err = err2
		return
	}
	return NumberExpression{lv - rv}, nil
}

func mul(left, right expression) (ne NumberExpression, err error) {
	lv, err2 := left.Value()
	if err2 != nil {
		err = err2
		return
	}
	rv, err2 := right.Value()
	if err2 != nil {
		err = err2
		return
	}
	return NumberExpression{lv * rv}, nil
}

func div(left, right expression) (ne NumberExpression, err error) {
	lv, err2 := left.Value()
	if err2 != nil {
		err = err2
		return
	}
	rv, err2 := right.Value()
	if err2 != nil {
		err = err2
		return
	}
	if rv == 0 {
		err = errors.New("Divide by zero!") // todo: refactor to a MathError type
		return
	}
	return NumberExpression{lv / rv}, nil
}

// Compute is the top level function of the interpreter.
func Compute(input string, vm variableMap) string {
	if vm == nil {
		vm = make(variableMap)
	}

	// Tokenization
	tokens, err := lexer.Lex(input)
	if err != nil {
		return err.Error()
	}

	// Parsing
	operatorStack, operantStack, err2 := parse(&tokenSliceEmitter{tokens: tokens}, vm)
	if err2 != nil {
		return err2.Error()
	}

	result, err3 := interpreting(operatorStack, operantStack)
	if err3 != nil {
		return err3.Error()
	}
	return strconv.FormatInt(result, 10)
}

type tokenSliceEmitter struct {
	tokens []lexer.Token
	i      int
}

func (emitter *tokenSliceEmitter) Next() lexer.Token {
	if emitter.i < len(emitter.tokens) {
		defer func() { emitter.i++ }()
		return emitter.tokens[emitter.i]
	}
	return nil
}

func interpreting(operatorStack []operator, operantStack []expression) (int64, error) {
	var operatorFunc operator
	var left, right expression
	for len(operatorStack) > 0 {
		operatorFunc, operatorStack = operatorStack[0], operatorStack[1:]
		left, operantStack = operantStack[0], operantStack[1:]
		right, operantStack = operantStack[0], operantStack[1:]
		result, err := operatorFunc(left, right)
		if err != nil {
			return 0, err
		}
		operantStack = append([]expression{result}, operantStack...)
	}
	return operantStack[0].Value()
}
