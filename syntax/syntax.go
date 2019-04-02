package syntax

import (
	"github.com/enabokov/language/lexis"
)

func Analyze(input lexis.TokenStream) TokenProgram {
	return program(input)
}
