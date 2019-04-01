package main

import (
	"fmt"
	"os"

	"github.com/enabokov/language/lexis"
	"github.com/enabokov/language/syntax"
)

func main() {
	filename := os.Args[1]
	tokenStream := lexis.Analyze(filename)
	fmt.Println(syntax.Analyze(tokenStream))
}
