package lexis

import (
	"strings"

	"github.com/enabokov/language/bnf"
)

type Token struct {
	Class string
	Value string
}

func getLexemes(sourceCode string) (lexemes []string) {
	return strings.Fields(sourceCode)
}

func getTokens(lexemes []string, bnfConfig bnf.BNF) []Token {
	var tokens []Token
	for _, lexeme := range lexemes {
		if checkValueInArray(lexeme, bnfConfig.Keywords) {
			tokens = append(tokens, Token{
				Class: `keyword`,
				Value: lexeme,
			})
		} else if checkValueInArray(lexeme, bnfConfig.PossibleType) {
			tokens = append(tokens, Token{
				Class: `type`,
				Value: lexeme,
			})
		} else if checkValueInArray(lexeme, bnfConfig.Punctuation) {
			tokens = append(tokens, Token{
				Class: `punctuation`,
				Value: lexeme,
			})
		} else if checkValueIsString(lexeme) {
			tokens = append(tokens, Token{
				Class: `string`,
				Value: lexeme,
			})
		} else if checkValueIsNumber(lexeme) {
			tokens = append(tokens, Token{
				Class: `number`,
				Value: lexeme,
			})
		} else if checkValueIsVariable(lexeme) {
			tokens = append(tokens, Token{
				Class: `variable`,
				Value: lexeme,
			})
		} else {
			tokens = append(tokens, Token{
				Class: `function`,
				Value: lexeme,
			})
		}
	}

	return tokens
}
