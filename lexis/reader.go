package lexis

import (
	"fmt"
	"log"
)

type Token struct {
	Class string
	Value string
}

const (
	ClassVariable    = "variable"
	ClassFunction    = "function"
	ClassKeyword     = "keyword"
	ClassCall        = "call"
	ClassString      = "string"
	ClassPunctuation = "punctuation"
	ClassOperator    = "operator"
	ClassNumber      = "number"
	ClassType        = "type"
	ClassBool        = "bool"
	ClassProgram     = "program"
)

func readWhile(input stream, predicate func(lexeme string) bool) (lexeme string) {
	for !input.eof() && predicate(input.peek()) {
		lexeme += input.next()
	}

	return lexeme
}

func readNumber(input stream) *Token {
	var hasDot = false
	number := readWhile(input,
		func(lexeme string) bool {
			if lexeme == "." {
				if hasDot {
					return false
				}
				hasDot = true
				return true
			}

			return isDigit(lexeme)
		},
	)

	return &Token{
		Class: ClassNumber,
		Value: number,
	}
}

func readIdentifier(input stream) *Token {
	id := readWhile(input, isIdentifier)

	var class = ClassVariable
	if isKeyword(id) {
		class = ClassKeyword
	}

	if isType(id) {
		class = ClassType
	}

	return &Token{
		Class: class,
		Value: id,
	}
}

func readCaller(input stream) *Token {
	caller := readWhile(input, isCall)

	return &Token{
		Class: ClassCall,
		Value: caller,
	}
}

func readEscaped(input stream, end string) string {
	var (
		escaped = false
		lexeme  string
	)

	input.next()
	for !input.eof() {
		ch := input.next()
		if escaped {
			lexeme += ch
			escaped = false
		} else if ch == "\\" {
			escaped = true
		} else if ch == end {
			break
		} else {
			lexeme += ch
		}
	}

	return lexeme
}

func readString(input stream) *Token {
	return &Token{
		Class: ClassString,
		Value: readEscaped(input, `"`),
	}
}

func skipComment(input stream) {
	readWhile(input,
		func(lexeme string) bool {
			return lexeme != "\n"
		},
	)

	input.next()
}

func readNext(input stream) (token *Token) {
	for {
		readWhile(input, isWhitespace)
		if input.eof() {
			return nil
		}

		ch := input.peek()
		if ch == `#` {
			skipComment(input)
			continue
		}

		if ch == `"` {
			return readString(input)
		}

		if isDigit(ch) {
			return readNumber(input)
		}

		if isIdentifierStart(ch) {
			return readIdentifier(input)
		}

		if isCall(ch) {
			return readCaller(input)
		}

		if isPunctuation(ch) {
			return &Token{
				Class: ClassPunctuation,
				Value: input.next(),
			}
		}

		if isOperator(ch) {
			return &Token{
				Class: ClassOperator,
				Value: readWhile(input, isOperator),
			}
		}

		log.Fatalln(input.croak(fmt.Sprintf("Can't handle: '%s'", ch)))
	}
}
