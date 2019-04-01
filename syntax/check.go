package syntax

import (
	"github.com/enabokov/language/lexis"
)

var priorities = map[string]int{
	`=`:  1,
	`||`: 2,
	`&&`: 3,
	`<`:  7, `>`: 7, `<=`: 7, `>=`: 7, `==`: 7, `!=`: 7,
	`+`: 10, `-`: 10,
	`*`: 20, `/`: 20, `%`: 20,
}

func isPackage(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassKeyword && token.Value == "package" {
		nextToken := input.Peek()
		if nextToken.Class == lexis.ClassVariable {
			return true
		}
	}

	return false
}

func isImport(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassKeyword && token.Value == "import" {
		nextToken := input.Peek()
		if nextToken.Class == lexis.ClassString {
			return true
		}
	}

	return false
}

func isFunction(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassKeyword && token.Value == "def" {
		nextToken := input.Peek()
		if nextToken.Class == lexis.ClassVariable {
			return true
		}
	}

	return false
}

func isVariable(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassKeyword && token.Value == "var" {
		nextToken := input.Peek()
		if nextToken.Class == lexis.ClassVariable {
			return true
		}
	}

	return false
}

func isCaller(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassVariable && input.Peek().Class == lexis.ClassCall {
		nextToken := input.Peek()
		if nextToken.Class == lexis.ClassCall {
			return true
		}
	}

	return false
}

func isAssignment(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassVariable && input.Peek().Class == lexis.ClassOperator && input.Peek().Value == `=` {
		return true
	}

	return false
}

func isCondition(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassKeyword && token.Value == "if" {
		nextToken := input.Peek()
		if nextToken.Class == lexis.ClassVariable {
			return true
		}
	}

	return false
}

func isBinaryExpression(input lexis.TokenStream, token *lexis.Token) bool {
	if token.Class == lexis.ClassVariable && input.Peek().Class == lexis.ClassOperator {
		return true
	}

	return false
}
