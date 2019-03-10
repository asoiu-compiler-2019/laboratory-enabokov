package syntax

import (
	"fmt"

	"github.com/enabokov/language/lexis"
)

func checkPackageExpression(tokens []lexis.Token) ([]lexis.Token, error) {
	var err error
	var currentToken lexis.Token
	var nextToken lexis.Token

	currentToken, tokens = tokens[0], tokens[1:]
	nextToken, tokens = tokens[0], tokens[1:]

	if currentToken.Class == "keyword" && currentToken.Value == "package" {
		if nextToken.Class == "variable" && nextToken.Value != "" {
			return tokens, nil
		}
	}

	err = fmt.Errorf(fmt.Sprintf("Failed to parse package expression: '%s %s'", currentToken.Value, nextToken.Value))
	return tokens, err
}

func _importCheck(current lexis.Token, next lexis.Token) (bool, error) {
	var err error
	if current.Class == "keyword" && current.Value == "import" {
		if next.Class == "string" && next.Value != "" {
			return true, nil
		}

		err = fmt.Errorf(fmt.Sprintf("Failed to parse import expression: '%s %s'. Add double qoutes", current.Value, next.Value))
	}

	return false, err
}

func checkImportExpression(tokens []lexis.Token) ([]lexis.Token, error) {
	var err error

	i := 0
	var ok bool
	for ok = true; ok; ok, err = _importCheck(tokens[i], tokens[i+1]) {
		tokens = tokens[2:]
		i = 0
	}

	return tokens, err
}
