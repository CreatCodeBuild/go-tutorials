package lexer

import (
	"text/scanner"
)

type ToekenReceiver struct {
	tokens []Token
}

func Lex(input string) ([]Token, *Error) {
	t := ToekenReceiver{}
	r := NewRuneEmitter(input)
	for next, err := StateBegin(r, &t); next != nil || err != nil; {
		if err != nil {
			return t.tokens, err
		}
		next, err = next(r, &t)
	}
	return t.tokens, nil
}

func (l *ToekenReceiver) Send(t Token) {
	l.tokens = append(l.tokens, t)
}

type RuneEmitter struct {
	ch   chan rune
	done chan struct{}
}

func NewRuneEmitter(input string) *RuneEmitter {
	r := RuneEmitter{ch: make(chan rune), done: make(chan struct{})}
	go func(r RuneEmitter) {
		for _, c := range input {
			r.ch <- c
		}
		r.ch <- scanner.EOF
		close(r.done)
	}(r)
	return &r
}

func (r *RuneEmitter) Next() (c rune) {
	select {
	case <-r.done:
		return scanner.EOF
	case next := <-r.ch:
		return next
	}
}
