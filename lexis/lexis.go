package lexis

import (
	"io/ioutil"
	"log"

	"github.com/enabokov/language/bnf"
)

var bnfConfig bnf.BNF

func init() {
	bnfConfig = bnf.Read()
}

func readFile(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	return string(file)
}

func Analyze(filename string) []Token {
	sourceCode := readFile(filename)
	lexemes := getLexemes(sourceCode)
	tokens := getTokens(lexemes, bnfConfig)
	return tokens
}
