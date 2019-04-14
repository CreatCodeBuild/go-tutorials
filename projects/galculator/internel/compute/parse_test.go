package compute

import (
	"fmt"
	"galculator/internel/lexer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsingIdentifierExpression(t *testing.T) {
	tokens, _ := lexer.Lex("a = 233")
	exp, err := parseIdentifierExpression(&tokenSliceEmitter{tokens: tokens[1:]}, lexer.Identifier{Value: "a"}, make(variableMap))
	require.NoError(t, err)
	if ae, ok := exp.(AssignmentExpression); ok {
		require.Equal(t, "a", ae.Name)
		value, err := ae.Value()
		require.Equal(t, int64(233), value)
		require.NoError(t, err)
	} else {
		a := require.New(t)
		a.FailNowf("not an assignment expression", "got %T", exp)
	}

	tokens, _ = lexer.Lex("a = 233 + 1")
	exp, err = parseIdentifierExpression(&tokenSliceEmitter{tokens: tokens[1:]}, lexer.Identifier{Value: "a"}, make(variableMap))
	require.NoError(t, err)
	if ae, ok := exp.(AssignmentExpression); ok {
		require.Equal(t, "a", ae.Name)
		value, err := ae.Value()
		require.Equal(t, int64(234), value)
		require.NoError(t, err)
	} else {
		a := require.New(t)
		a.FailNowf("not an assignment expression", "got %T", exp)
	}

	tokens, _ = lexer.Lex("a = 1 + 2 + 3")
	exp, err = parseIdentifierExpression(&tokenSliceEmitter{tokens: tokens[1:]}, lexer.Identifier{Value: "a"}, make(variableMap))
	require.NoError(t, err)
	if ae, ok := exp.(AssignmentExpression); ok {
		require.Equal(t, "a", ae.Name)
		value, err := ae.Value()
		require.Equal(t, int64(6), value)
		require.NoError(t, err)
	} else {
		a := require.New(t)
		a.FailNowf("not an assignment expression", "got %T", exp)
	}

	tokens, _ = lexer.Lex("a = 1 + (2 + 3)")
	exp, err = parseIdentifierExpression(&tokenSliceEmitter{tokens: tokens[1:]}, lexer.Identifier{Value: "a"}, make(variableMap))
	require.NoError(t, err)
	if ae, ok := exp.(AssignmentExpression); ok {
		require.Equal(t, "a", ae.Name)
		value, err := ae.Value()
		require.Equal(t, int64(6), value)
		require.NoError(t, err)
	} else {
		a := require.New(t)
		a.FailNowf("not an assignment expression", "got %T", exp)
	}

	tokens, _ = lexer.Lex("a = a")
	exp, err = parseIdentifierExpression(&tokenSliceEmitter{tokens: tokens[1:]}, lexer.Identifier{Value: "a"}, make(variableMap))
	require.NoError(t, err)
	value, err := exp.Value()
	require.EqualError(t, err, "Variable a is not defined")
	require.Equal(t, int64(0), value)

	tokens, _ = lexer.Lex("a = (1 + 2 + a) + 3")
	exp, err = parseIdentifierExpression(&tokenSliceEmitter{tokens: tokens[1:]}, lexer.Identifier{Value: "a"}, make(variableMap))
	require.NoError(t, err)
	if ae, ok := exp.(AssignmentExpression); ok {
		require.Equal(t, "a", ae.Name)
		value, err := ae.Value()
		require.EqualError(t, err, "Variable a is not defined")
		require.Equal(t, int64(0), value)
	} else {
		a := require.New(t)
		a.FailNowf("not an assignment expression", "got %T", exp)
	}
}

func TestParsingParenthesisExpression(t *testing.T) {
	tokens, _ := lexer.Lex("(1 + 2)")
	exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
	require.NoError(t, err)
	require.IsType(t, ParenthesesExpression{}, exp)
	result, err := exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(3), result)

	tokens, _ = lexer.Lex("(233)")
	exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
	require.NoError(t, err)
	require.IsType(t, ParenthesesExpression{}, exp)
	result, err = exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(233), result)

	tokens, _ = lexer.Lex("(1 + (2))")
	exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
	require.NoError(t, err)
	require.IsType(t, ParenthesesExpression{}, exp)
	result, err = exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(3), result)

	tokens, _ = lexer.Lex("1 + 2")
	exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
	require.NoError(t, err)
	require.IsType(t, OperatorExpression{}, exp)
	result, err = exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(3), result)

	tokens, _ = lexer.Lex("((1) + 2)")
	exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
	require.NoError(t, err)
	require.IsType(t, ParenthesesExpression{}, exp)
	result, err = exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(3), result)

	tokens, _ = lexer.Lex("(((100))+233+(1-100))")
	exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
	require.NoError(t, err)
	require.IsType(t, ParenthesesExpression{}, exp)
	result, err = exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(234), result)
}

func TestParsingExpression(t *testing.T) {

	t.Run("parse negative number", func(t2 *testing.T) {
		tokens, _ := lexer.Lex("-1-1")
		exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		if opExp, ok := exp.(OperatorExpression); ok {
			fmt.Println(opExp.Left, opExp.Op, opExp.Right)
		}

		result, err := exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(-1), result)
	})

	t.Run("parse identifier exp", func(t2 *testing.T) {
		tokens, _ := lexer.Lex("(a)")
		exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, ParenthesesExpression{}, exp)
		result, err := exp.Value()
		require.EqualError(t, err, "Variable a is not defined")
		require.Equal(t, int64(0), result)
	})

	t.Run("parse +-*/ exp", func(t2 *testing.T) {
		tokens, _ := lexer.Lex("1+1")
		exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err := exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(2), result)

		tokens, _ = lexer.Lex("1-2")
		exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err = exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(-1), result)

		tokens, _ = lexer.Lex("666*998")
		exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err = exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(666*998), result)

		tokens, _ = lexer.Lex("6/2")
		exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err = exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(3), result)

		tokens, _ = lexer.Lex("6/0")
		exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err = exp.Value()
		require.EqualError(t, err, "Divide by zero!")
		require.Equal(t, int64(0), result)

		tokens, _ = lexer.Lex("1+6/1")
		exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err = exp.Value()
		require.EqualError(t, err, "Divide by zero!")
		require.Equal(t, int64(0), result)
	})

	t.Run("parse number exp", func(t2 *testing.T) {
		tokens, _ := lexer.Lex("233")
		exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, NumberExpression{}, exp)

		result, err := exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(233), result)
	})

	t.Run("parse = exp", func(t2 *testing.T) {
		tokens, _ := lexer.Lex("a = 1 + 2 + 3")

		exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, make(variableMap))
		require.NoError(t, err)
		require.IsType(t, AssignmentExpression{}, exp)

		result, err := exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(6), result)
	})

	t.Run("signed number", func(t2 *testing.T) {
		tokens, _ := lexer.Lex("+1")
		exp, err := parseExpression(&tokenSliceEmitter{tokens: tokens}, nil)
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err := exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(1), result)

		tokens, _ = lexer.Lex("-1")
		exp, err = parseExpression(&tokenSliceEmitter{tokens: tokens}, nil)
		require.NoError(t, err)
		require.IsType(t, OperatorExpression{}, exp)
		result, err = exp.Value()
		require.NoError(t, err)
		require.Equal(t, int64(-1), result)
	})

}
