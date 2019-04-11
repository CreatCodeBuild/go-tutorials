package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	tokens, _ := Lex("1+1")
	require.Equal(t,
		[]Token{
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
		},
		tokens)

	tokens, _ = Lex("1 + 1")
	require.Equal(t,
		[]Token{
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
		},
		tokens)

	tokens, _ = Lex("1 + 1 0000 0000") // support white space between numbers for better big number literal displaying
	require.Equal(t,
		[]Token{
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "100000000"},
		},
		tokens)

	tokens, _ = Lex("1 + 1 - 20") // support multiple digit numbers
	require.Equal(t,
		[]Token{
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
			Operator{Value: "-"},
			Number{Value: "20"},
		},
		tokens)
}

func TestLexerWrongSyntaxValidToken(t *testing.T) {
	// Tests in this function all have incorrect syntax but they are not abled to be caught by a finite state machine.
	// They can only be caught at the parsing level by a context free language.

	tokens, _ := Lex("((2)")
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			LeftParentheses{},
			Number{Value: "2"},
			RightParentheses{},
		},
		tokens)
}

func TestLexerParenthese(t *testing.T) {
	tokens, _ := Lex("(1+1")
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
		},
		tokens)

	tokens, _ = Lex("(1+1)")
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
			RightParentheses{},
		},
		tokens)

	tokens, _ = Lex("(1+1)/99") // support multiple digit numbers
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
			RightParentheses{},
			Operator{Value: "/"},
			Number{Value: "99"},
		},
		tokens)

	tokens, _ = Lex("(1+(1+1))")
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			Number{Value: "1"},
			Operator{Value: "+"},
			LeftParentheses{},
			Number{Value: "1"},
			Operator{Value: "+"},
			Number{Value: "1"},
			RightParentheses{},
			RightParentheses{},
		},
		tokens)
}

func TestLexerError(t *testing.T) {
	_, err := Lex("1 + function")
	require.EqualError(t, err, "Lexing Error: character 'f' is not expected after an operator")

	_, err = Lex("1 + +")
	require.EqualError(t, err, "Lexing Error: character '+' is not expected after an operator")

	_, err = Lex("1 function")
	require.EqualError(t, err, "Lexing Error: character 'f' is not expected after a number")

	_, err = Lex("1+")
	require.EqualError(t, err, "Lexing Error: An operator must be followed by an expression")

	_, err = Lex("(")
	require.EqualError(t, err, "Lexing Error: A ( must be followed by an expression")

	_, err = Lex("()")
	require.EqualError(t, err, "Lexing Error: An expression must be in between ( and )")

	_, err = Lex("(()")
	require.EqualError(t, err, "Lexing Error: An expression must be in between ( and )")
}

func TestAssignment(t *testing.T) {
	tokens, err := Lex("babcccccc = 2")
	require.Nil(t, err)
	require.Equal(t,
		[]Token{
			Identifier{Value: "babcccccc"},
			EqualSign{},
			Number{Value: "2"},
		},
		tokens)

	tokens, err = Lex("1 +a")
	require.Nil(t, err)
	require.Equal(t,
		[]Token{
			Number{Value: "1"},
			Operator{Value: "+"},
			Identifier{Value: "a"},
		},
		tokens)

	tokens, err = Lex("(a)")
	require.Nil(t, err)
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			Identifier{Value: "a"},
			RightParentheses{},
		},
		tokens)

	tokens, err = Lex("(1)a")
	require.EqualError(t, err, "Lexing Error: character 'a' is not expected after )")
	require.Equal(t,
		[]Token{
			LeftParentheses{},
			Number{Value: "1"},
			RightParentheses{},
		},
		tokens)

	tokens, err = Lex("a")
	require.Nil(t, err)
	require.Equal(t,
		[]Token{
			Identifier{Value: "a"},
		},
		tokens)

	tokens, err = Lex("babcccccc == 2")
	require.Nil(t, err)
	require.Equal(t,
		[]Token{
			Identifier{Value: "babcccccc"},
			EqualSign{},
			EqualSign{},
			Number{Value: "2"},
		},
		tokens)
}
