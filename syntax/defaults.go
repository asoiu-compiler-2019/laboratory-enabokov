package syntax

type astNode interface{}

type tokenPrimitive struct {
	Class string
	Value string
}

type tokenVariable struct {
	Class string
	Name  string
	Type  string
}

type tokenPackage struct {
	Class string
	Value string
}

type tokenImport struct {
	Class string
	Value string
}

type tokenCondition struct {
	Class     string
	Condition tokenBinaryExprOrAssign
	Do        []astNode
	Else      astNode
}

type tokenCall struct {
	Class string
	Func  tokenVariable
	Args  []tokenVariable
}

type tokenBinaryExprOrAssign struct {
	Class    string
	Operator string
	Left     tokenPrimitive
	Right    *tokenBinaryExprOrAssign
}

type tokenFunction struct {
	Class  string
	Name   string
	Params []tokenVariable
	Body   []astNode
}

type tokenProgram struct {
	Class      string
	Expression []astNode
}
