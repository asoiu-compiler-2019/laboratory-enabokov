package semantics

import (
	"github.com/enabokov/language/syntax"
)

func Analyze(ast syntax.TokenProgram) error {
	// fmt.Printf("%# v", pretty.Formatter(ast))
	return scan(ast)
}
