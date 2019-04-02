package syntax

type ASTNode interface{}

type TokenPrimitive struct {
	Class string
	Value string
}

type TokenVariable struct {
	Class string
	Name  string
	Type  string
}

type TokenPackage struct {
	Class string
	Value string
}

type TokenImport struct {
	Class string
	Value string
}

type TokenCondition struct {
	Class     string
	Condition TokenBinaryOrAssign
	Do        []ASTNode
	Else      ASTNode
}

type TokenCall struct {
	Class string
	Func  TokenVariable
	Args  []TokenVariable
}

type TokenBinaryOrAssign struct {
	Class    string
	Operator string
	Left     TokenPrimitive
	Right    *TokenBinaryOrAssign
}

type TokenFunction struct {
	Class  string
	Name   string
	Params []TokenVariable
	Body   []ASTNode
}

type TokenProgram struct {
	Class      string
	Expression []ASTNode
}
