package compute

import (
	"fmt"
	"galculator/internel/lexer"
	"strconv"
)

type operator func(left expression, right expression) NumberExpression

func add(left, right expression) NumberExpression {
	return NumberExpression{number: left.Value() + right.Value()}
}

func sub(left, right expression) NumberExpression {
	return NumberExpression{number: left.Value() - right.Value()}
}

func mul(left, right expression) NumberExpression {
	return NumberExpression{number: left.Value() * right.Value()}
}

func div(left, right expression) NumberExpression {
	return NumberExpression{number: left.Value() / right.Value()}
}

func Compute(input string) string {
	// Tokenization
	tokens, err := lexer.Lex(input)
	if err != nil {
		return err.Error()
	}

	// Parsing
	operatorStack, operantStack, err2 := parse(tokens)
	if err2 != nil {
		return err2.Error()
	}

	return strconv.FormatInt(interpreting(operatorStack, operantStack), 10)
}

func interpreting(operatorStack []operator, operantStack []expression) int64 {
	var operatorFunc operator
	var left, right expression
	for len(operatorStack) > 0 {
		operatorFunc, operatorStack = operatorStack[0], operatorStack[1:]
		left, operantStack = operantStack[0], operantStack[1:]
		right, operantStack = operantStack[0], operantStack[1:]
		operantStack = append([]expression{operatorFunc(left, right)}, operantStack...)
	}
	return operantStack[0].Value()
}

func parse(tokens []lexer.Token) (operatorStack []operator, operantStack []expression, err error) {
	operators := map[lexer.Operator]operator{
		lexer.Add: add,
		lexer.Sub: sub,
		lexer.Mul: mul,
		lexer.Div: div,
	}

	for i := 0; i < len(tokens); {
		inc := 1
		switch t := tokens[i].(type) {
		case lexer.Operator:
			if operator, ok := operators[t]; ok {
				operatorStack = append(operatorStack, operator)
			} else {
				panic("?")
			}
		case lexer.Number:
			integer, err := strconv.ParseInt(t.Value, 10, 64)
			if err != nil {
				panic(err)
			}
			operantStack = append(operantStack, NumberExpression{number: integer})
		case lexer.LeftParentheses:
			pe, read, err2 := ParseParenthesisExpression(tokens[i+1:])
			if err2 != nil {
				err = err2
				return
			}
			operantStack = append(operantStack, pe)
			inc = read + 1 // todo: coordinate the current read position of tokens...
		default:
			fmt.Println("Warning:", t.Type(), t.Literal(), "is ignored")
		}
		i += inc
	}
	return
}

// ParseParenthesisExpression parse a sequence of tokens and return the first full parenthesis expression.
// ( 1 + 1 ) + 1 ) will return expression"((1+1)+1)"
// This function always inserts a leading ( at the beginning of the tokens.
func ParseParenthesisExpression(s []lexer.Token) (ParenthesesExpression, int, error) {
	count := 1
	for i, token := range s {
		switch token.(type) {
		case lexer.LeftParentheses:
			count++
		case lexer.RightParentheses:
			count--
		}
		if count == 0 {
			operatorStack, operantStack, err := parse(s[:i])
			return ParenthesesExpression{
				OperatorStack: operatorStack,
				OperantStack:  operantStack,
			}, i + 1, err
		}
	}
	return ParenthesesExpression{}, len(s), &ParsingError{fmt.Sprintf("Missing %d ) parentheses", count)}
}

type expression interface {
	Value() int64
}

type NumberExpression struct {
	number int64
}

func (ne NumberExpression) Value() int64 {
	return ne.number
}

type ParenthesesExpression struct {
	OperatorStack []operator
	OperantStack  []expression
}

func (pe ParenthesesExpression) Value() int64 {
	return interpreting(pe.OperatorStack, pe.OperantStack)
}
