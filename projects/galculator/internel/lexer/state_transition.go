package lexer

import "bytes"

func validIdentifierRune(r rune) bool {
	return bytes.ContainsRune([]byte("abc"), r)
}

func validOperator(r rune) bool {
	return bytes.ContainsRune([]byte("+-*/"), r)
}
