package lexer

import (
	"fmt"
	"text/scanner"
)

// This interface happens to match https://golang.org/pkg/text/scanner/#Scanner.Next
type runeEmitter interface {
	Next() rune
}

type tokenReciver interface {
	Send(Token)
}

type StateFunc func(runeEmitter, tokenReciver) (StateFunc, *Error)

func StateBegin(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {
	next := r.Next()
	switch next {
	case ' ': // Single quote in Go is for character.
		return StateBegin, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return StateNumber([]rune{next}), nil
	case '(':
		tokens.Send(LeftParentheses{})
		return StateLeftParenthesis, nil
	case '=':
		tokens.Send(EqualSign{})
		return StateBegin, nil
	default:
		if validIdentifierRune(next) {
			return StateIdentifier([]rune{next}), nil
		}
		if validOperator(next) {
			tokens.Send(Operator{Value: string(next)})
			return StateOperator, nil
		}
		if next == scanner.EOF {
			return nil, &Error{"StateBegin reaches EOF."}
		}
		return nil, &Error{fmt.Sprintf("character '%s' is not expected", string(next))}
	}
}

func StateLeftParenthesis(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {
	next := r.Next()
	switch next {
	case ' ': // Single quote in Go is for character.
		return StateLeftParenthesis, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return StateNumber([]rune{next}), nil
	case '+', '-', '*', '/':
		tokens.Send(Operator{Value: string(next)})
		return StateOperator, nil
	case '(':
		tokens.Send(LeftParentheses{})
		return StateLeftParenthesis, nil
	case ')':
		return nil, &Error{"An expression must be in between ( and )"}
	case scanner.EOF:
		return nil, &Error{"A ( must be followed by an expression"}
	default:
		if validIdentifierRune(next) {
			return StateIdentifier([]rune{next}), nil
		}
		return nil, &Error{fmt.Sprintf("character '%s' is not expected", string(next))}
	}
}

func StateRightParentheses(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {
	next := r.Next()
	switch next {
	case ' ': // Single quote in Go is for character.
		return StateRightParentheses, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return nil, &Error{fmt.Sprintf("number is not expected after )")}
	case '+', '-', '*', '/':
		tokens.Send(Operator{Value: string(next)})
		return StateOperator, nil
	case '(':
		return nil, &Error{fmt.Sprintf("( is not expected after )")}
	case ')':
		tokens.Send(RightParentheses{})
		return StateRightParentheses, nil
	case scanner.EOF:
		return nil, nil
	default:
		return nil, &Error{fmt.Sprintf("character '%s' is not expected after )", string(next))}
	}
}

func StateNumber(read []rune) StateFunc {

	return func(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {
		next := r.Next()
		switch next {
		case ' ':
			return StateNumber(read), nil
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return StateNumber(append(read, next)), nil
		case '+', '-', '*', '/':
			tokens.Send(Number{
				Value: string(read),
			})
			tokens.Send(Operator{
				Value: string(next),
			})
			return StateOperator, nil
		case ')':
			tokens.Send(Number{
				Value: string(read),
			})
			tokens.Send(RightParentheses{})
			return StateRightParentheses, nil
		case scanner.EOF:
			tokens.Send(Number{Value: string(read)})
			return nil, nil
		default:
			return nil, &Error{fmt.Sprintf("character '%s' is not expected after a number", string(next))}
		}
	}
}

// // StateBeginOperatorOrNumber : when read an operator, it could be a number with a sign or an operator by its own.
// func StateBeginOperatorOrNumber(operator rune) StateFunc {
// 	return func(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {

// 		if operator == '*' || operator == '/' {
// 			tokens.Send(Operator{Value: string(operator)})
// 			return StateOperator, nil
// 		}

// 		next := r.Next()
// 		switch next {
// 		case ' ':
// 			return StateOperator, nil
// 		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
// 			return StateNumber([]rune{next}), nil
// 		case '(':
// 			tokens.Send(LeftParentheses{})
// 			return StateLeftParenthesis, nil
// 		case scanner.EOF:
// 			return nil, &Error{"An operator must be followed by an expression"}
// 		default:
// 			if validIdentifierRune(next) {
// 				return StateIdentifier([]rune{next}), nil
// 			}
// 			return nil, &Error{fmt.Sprintf("character '%s' is not expected after an operator", string(next))}
// 		}
// 	}
// }

func StateOperator(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {
	next := r.Next()
	switch next {
	case ' ':
		return StateOperator, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return StateNumber([]rune{next}), nil
	case '(':
		tokens.Send(LeftParentheses{})
		return StateLeftParenthesis, nil
	case scanner.EOF:
		return nil, &Error{"An operator must be followed by an expression"}
	default:
		if validIdentifierRune(next) {
			return StateIdentifier([]rune{next}), nil
		}
		return nil, &Error{fmt.Sprintf("character '%s' is not expected after an operator", string(next))}
	}
}

// StateIdentifier tokenize a variable name.
// English letter + Chinese
func StateIdentifier(read []rune) StateFunc {
	return func(r runeEmitter, tokens tokenReciver) (StateFunc, *Error) {
		next := r.Next()
		switch next {
		case scanner.EOF:
			tokens.Send(Identifier{Value: string(read)})
			return nil, nil
		case ')':
			tokens.Send(Identifier{Value: string(read)})
			tokens.Send(RightParentheses{})
			return StateRightParentheses, nil
		case '=':
			tokens.Send(Identifier{Value: string(read)})
			tokens.Send(EqualSign{})
			return StateBegin, nil
		case '+', '-', '*', '/':
			tokens.Send(Identifier{Value: string(read)})
			tokens.Send(Operator{Value: string(next)})
			return StateBegin, nil
		default:
			if validIdentifierRune(next) {
				return StateIdentifier(append(read, next)), nil
			}
			tokens.Send(Identifier{Value: string(read)})
			return StateBegin, nil
		}
	}
}
