package main

import (
	"fmt"
	"os"

	"github.com/enabokov/language/lexis"
	"github.com/enabokov/language/syntax"
)

func main() {
	filename := os.Args[1]
	tokens := lexis.Analyze(filename)
	ok := syntax.Analyze(tokens)
	fmt.Println(ok)
}
