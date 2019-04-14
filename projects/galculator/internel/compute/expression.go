package compute

import (
	"errors"
	"fmt"
	"galculator/internel/lexer"
)

type expression interface {
	Value() (int64, error)
}

// NumberExpression represents an integer value
type NumberExpression struct {
	number int64
}

func (ne NumberExpression) Value() (int64, error) {
	return ne.number, nil
}

// OperatorExpression is
// exp op exp
// If op is + or -, the first exp can be nil.
// which means, nil + 1 is positive 1, nil - 1 is negative 1.
type OperatorExpression struct {
	Op    lexer.Operator
	Left  expression
	Right expression
}

func (oe OperatorExpression) Value() (int64, error) {

	var lv int64
	if oe.Left == nil && (oe.Op == lexer.Add || oe.Op == lexer.Sub) {
		lv = 0
	} else {
		var err error
		lv, err = oe.Left.Value()
		if err != nil {
			return 0, err
		}
	}

	rv, err2 := oe.Right.Value()
	if err2 != nil {
		return 0, err2
	}
	switch oe.Op {
	case lexer.Add:
		return lv + rv, nil
	case lexer.Sub:
		return lv - rv, nil
	case lexer.Mul:
		return lv * rv, nil
	case lexer.Div:
		if rv == 0 {
			return 0, errors.New("Divide by zero!") // todo: refactor to a MathError type
		}
		return lv / rv, nil
	default:
		return 0, errors.New("Op not recognized")
	}
}

// ParenthesesExpression is
// ( expression )
type ParenthesesExpression struct {
	OperatorStack []operator
	OperantStack  []expression
}

func (pe ParenthesesExpression) Value() (int64, error) {
	return interpreting(pe.OperatorStack, pe.OperantStack)
}

// AssignmentExpression is
// identifier = expression
type AssignmentExpression struct {
	Name       string
	expression expression
	vm         variableMap
}

func (ae AssignmentExpression) Value() (int64, error) {
	v, err := ae.expression.Value()
	if err != nil {
		return 0, err
	}
	ae.vm[ae.Name] = v
	return v, nil
}

// IdentifierExpression is an identifier by itself.
type IdentifierExpression struct {
	Name     string
	ValueMap map[string]int64
}

func (ie IdentifierExpression) Value() (int64, error) {
	if v, ok := ie.ValueMap[ie.Name]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("Variable %s is not defined", ie.Name)
}
