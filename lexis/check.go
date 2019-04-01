package lexis

import (
	"regexp"
	"strconv"
)

func isKeyword(lexeme string) bool {
	for _, key := range BnfConfig.Keywords {
		if key == lexeme {
			return true
		}
	}

	return false
}

func isType(lexeme string) bool {
	for _, key := range BnfConfig.Types {
		if key == lexeme {
			return true
		}
	}

	return false
}

func isDigit(lexeme string) bool {
	if _, err := strconv.Atoi(lexeme); err != nil {
		return false
	}

	return true
}

func isIdentifierStart(lexeme string) bool {
	return regexp.MustCompile(`[a-zA-Z_]`).MatchString(lexeme)
}

func isIdentifier(lexeme string) bool {
	return isIdentifierStart(lexeme) || regexp.MustCompile(`[0-9-]`).MatchString(lexeme)
}

func isOperator(lexeme string) bool {
	for _, operator := range BnfConfig.Operators {
		if operator == lexeme {
			return true
		}
	}

	return false
}

func isPunctuation(lexeme string) bool {
	for _, key := range BnfConfig.Punctuation {
		if key == lexeme {
			return true
		}
	}

	return false
}

func isWhitespace(lexeme string) bool {
	return regexp.MustCompile(`[[:space:]]`).MatchString(lexeme)
}

func isCall(lexeme string) bool {
	return regexp.MustCompile(`[a-zA-Z_0-9.]`).MatchString(lexeme)
}
