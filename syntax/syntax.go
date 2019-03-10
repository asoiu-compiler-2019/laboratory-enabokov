package syntax

import (
	"log"

	"github.com/enabokov/language/lexis"
)

func Analyze(tokens []lexis.Token) bool {
	tokens, err := checkPackageExpression(tokens)
	if err != nil {
		log.Fatalln(err)
	}

	tokens, err = checkImportExpression(tokens)
	if err != nil {
		log.Fatalln(err)
	}

	return true
}
