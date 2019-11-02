package lexer

import "errors"

// Token is ..
type Token struct {
	Type int
	Literal []rune
}
// var Number = 0

// var num := Token{
// 	Type: Number
// 	Literal: "123"
// }


func Lex(input string) ([]Token, error) {

	const (
		START int = iota
		NUMBER
		OPERATOR
	)

	state := START
	var read []rune
	var tokens []Token
	for _, letter := range input {
		switch state {
		case START:
			if !isDigit(letter) {
				return nil, errors.New("only expect digit at the beginning")
			}
			read = append(read, letter)
			state = NUMBER
		case NUMBER:
			if isDigit(letter) {
				read = append(read, letter)
				state = NUMBER
			} else if isOperator(letter) {
				tokens = append(tokens, Token{
					Type: NUMBER,
					Literal: read,
				})
				read = nil
				state = OPERATOR
			} else if letter == ' ' {
				continue
			} else {
				return nil, errors.New("error")
			}
		case OPERATOR:
		}
	}
}
