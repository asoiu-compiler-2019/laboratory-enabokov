package syntax

import (
	"github.com/enabokov/language/lexis"
)

func Analyze(input lexis.TokenStream) bool {
	return program(input)
}
