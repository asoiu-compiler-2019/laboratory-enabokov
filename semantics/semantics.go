package semantics

import (
	"fmt"

	"github.com/enabokov/language/syntax"
	"github.com/kr/pretty"
)

func Analyze(ast syntax.TokenProgram) error {
	fmt.Printf("%# v", pretty.Formatter(ast))
	return scan(ast)
}
