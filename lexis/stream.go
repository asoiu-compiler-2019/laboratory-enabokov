package lexis

import (
	"fmt"
)

type stream struct {
	next  func() string
	peek  func() string
	eof   func() bool
	croak func(string) error
}

func readInputStream(input []string) stream {
	var (
		line = 0
		col  = 0
	)

	next := func() (ch string) {
		if line < len(input) && col < len(input[line]) {
			ch = string(input[line][col])
			col++
		}

		if ch == "\n" {
			line++
			col = 0
		}

		return ch
	}

	peek := func() (ch string) {
		if line < len(input) && col < len(input[line]) {
			ch = string(input[line][col])
		}

		return ch
	}

	eof := func() bool {
		return peek() == ""
	}

	croak := func(msg string) error {
		return fmt.Errorf("(Line: %d, Column: %d) %s", line+1, col+1, msg)
	}

	return stream{next, peek, eof, croak}
}

type tokenStream struct {
	next  func() *Token
	peek  func() *Token
	eof   func() bool
	croak func(string) error
}

func readTokenStream(input stream) tokenStream {
	var current *Token

	next := func() (token *Token) {
		token = current
		current = nil

		if token == nil {
			return readNext(input)
		}

		return token
	}

	peek := func() *Token {
		if current == nil {
			current = readNext(input)
		}

		return current
	}

	eof := func() bool {
		return peek() == nil
	}

	return tokenStream{next, peek, eof, input.croak}
}
