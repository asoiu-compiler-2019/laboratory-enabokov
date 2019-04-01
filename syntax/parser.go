package syntax

import (
	"fmt"
	"log"
	"strconv"

	"github.com/enabokov/language/lexis"
	"github.com/kr/pretty"
)

func parseParams(input lexis.TokenStream, params *[]tokenVariable, requiredParam bool) error {
	nextToken := input.Peek()
	if nextToken.Class == lexis.ClassVariable {
		if !requiredParam {
			return input.Croak(fmt.Sprintf("Got `%s`. Expected `,` or `)`", input.Peek().Value))
		}

		t := input.Next()

		*params = append(
			*params,
			tokenVariable{
				Class: t.Class,
				Name:  t.Value,
			},
		)
		return parseParams(input, params, false)
	}

	if nextToken.Class == lexis.ClassPunctuation && nextToken.Value == "," {
		if requiredParam {
			return input.Croak(fmt.Sprintf("Got `%s`. Expected param", input.Peek().Value))
		}

		input.Next()
		return parseParams(input, params, true)
	}

	if requiredParam {
		return input.Croak(fmt.Sprintf("Got `%s`. Expected param", input.Peek().Value))
	}

	if input.Peek().Class == lexis.ClassPunctuation && input.Peek().Value == `)` {
		input.Next()
	}

	return nil
}

func parseArgs(input lexis.TokenStream, params *[]tokenVariable, requiredParam bool) error {
	nextToken := input.Peek()
	if nextToken.Class == lexis.ClassVariable {
		if !requiredParam {
			return input.Croak(fmt.Sprintf("Got `%s`. Expected `,` or `)`", input.Peek().Value))
		}

		t := input.Next()
		*params = append(
			*params,
			tokenVariable{
				Class: t.Class,
				Name:  t.Value,
			},
		)
		return parseParams(input, params, false)
	}

	if nextToken.Class == lexis.ClassPunctuation && nextToken.Value == "," {
		if requiredParam {
			return input.Croak(fmt.Sprintf("Got `%s`. Expected param", input.Peek().Value))
		}

		input.Next()
		return parseParams(input, params, true)
	}

	if requiredParam {
		return input.Croak(fmt.Sprintf("Got `%s`. Expected param", input.Peek().Value))
	}

	if input.Peek().Class == lexis.ClassPunctuation && input.Peek().Value == `)` {
		input.Next()
	}

	return nil
}

func parseBody(input lexis.TokenStream, body *[]astNode) error {
	token := input.Next()
	if token.Class != lexis.ClassPunctuation || token.Value != `{` {
		return input.Croak(fmt.Sprintf("Got `%s`. Expected {", token.Value))
	}

	for input.Peek().Class != lexis.ClassPunctuation || input.Peek().Value != `}` {
		astNode, err := expression(input)
		if err != nil {
			return input.Croak("Oooops" + err.Error())
		}

		*body = append(*body, astNode)
	}

	input.Next()
	return nil
}

func parseFunction(input lexis.TokenStream) (tokenFunction, error) {
	// get function name
	name := input.Next().Value

	var (
		params []tokenVariable
		err    error
	)

	if token := input.Next(); token.Class != lexis.ClassPunctuation || token.Value != `(` {
		return tokenFunction{}, input.Croak(fmt.Sprintf("Got `%s`. Expected `(`", token.Value))
	}

	if err = parseParams(input, &params, true); err != nil {
		return tokenFunction{}, err
	}

	var body []astNode
	err = parseBody(input, &body)

	return tokenFunction{
		Class:  `function`,
		Name:   name,
		Params: params,
		Body:   body,
	}, err
}

func parsePackage(input lexis.TokenStream) (tokenPackage, error) {
	return tokenPackage{Class: `package`, Value: input.Next().Value}, nil
}

func parseVariable(input lexis.TokenStream) (tokenVariable, error) {
	name := input.Next().Value
	token := input.Next()
	if token.Class != lexis.ClassType {
		return tokenVariable{}, input.Croak(fmt.Sprintf("Got `%s`. Expected type of variable", token.Value))
	}

	return tokenVariable{Class: `variable`, Name: name, Type: token.Value}, nil
}

func parseImport(input lexis.TokenStream) (tokenImport, error) {
	return tokenImport{Class: `import`, Value: input.Next().Value}, nil
}

func parseCaller(input lexis.TokenStream, token *lexis.Token) (tokenCall, error) {
	nextToken := input.Next()
	if input.Peek().Class == lexis.ClassPunctuation && input.Peek().Value == "(" {
		input.Next()
	}

	var args []tokenVariable

	err := parseArgs(input, &args, true)
	return tokenCall{
		Class: `caller`,
		Func:  tokenVariable{Class: token.Class, Name: token.Value + nextToken.Value},
		Args:  args,
	}, err
}

func _parse(operators []string, numbers []string, result *tokenBinaryExprOrAssign) error {
	var op string
	var num string

	if len(operators) == 0 && len(numbers) == 0 {
		return nil
	}

	if len(operators) > 0 {
		op, operators = operators[0], operators[1:]
		result.Operator = op
	}

	if len(numbers) > 0 {
		num, numbers = numbers[0], numbers[1:]
		result.Left = tokenPrimitive{
			Class: `number`,
			Value: num,
		}
	}

	result.Class = `binary`
	result.Right = &tokenBinaryExprOrAssign{}
	return _parse(operators, numbers, result.Right)
}

func parseBinaryExpression(input lexis.TokenStream, token *lexis.Token) (tokenBinaryExprOrAssign, error) {
	var numbers []string
	var operators []string

	for {
		token := input.Peek()

		if token.Class == lexis.ClassOperator {
			operators = append(operators, token.Value)
			input.Next()
			continue
		}

		if token.Class == lexis.ClassPunctuation && (token.Value == `(` || token.Value == `)`) {
			operators = append(operators, token.Value)
			input.Next()
			continue
		}

		if token.Class == lexis.ClassNumber {
			numbers = append(numbers, token.Value)
			input.Next()
			continue
		}

		break
	}

	var res = tokenBinaryExprOrAssign{Class: `assignment`}
	err := _parse(operators, numbers, &res)
	return res, err
}

func parseAssignment(input lexis.TokenStream, token *lexis.Token) (tokenBinaryExprOrAssign, error) {
	nextToken := input.Next()
	number := input.Peek()

	res, err := parseBinaryExpression(input, number)

	if _, err := strconv.Atoi(token.Value); err == nil {
		return tokenBinaryExprOrAssign{}, input.Croak(fmt.Sprintf("Got %s. Cannot assign to number", token.Value))
	}

	return tokenBinaryExprOrAssign{
		Class:    `assignment`,
		Operator: nextToken.Value,
		Left:     tokenPrimitive{Class: `variable`, Value: token.Value},
		Right:    &res,
	}, err
}

func parseCondition(input lexis.TokenStream, token *lexis.Token) (tokenCondition, error) {
	input.Next()

	var body []astNode
	cond, err := parseBinaryExpression(input, token)
	err = parseBody(input, &body)

	return tokenCondition{
		Class:     `condition`,
		Condition: cond,
		Do:        body,
	}, err
}

func expression(input lexis.TokenStream) (astNode, error) {
	token := input.Next()
	switch {
	default:
		return nil, input.Croak(fmt.Sprintf("Failed to parse `%s`", token.Value))
	case isPackage(input, token):
		return parsePackage(input)
	case isImport(input, token):
		return parseImport(input)
	case isFunction(input, token):
		return parseFunction(input)
	case isVariable(input, token):
		return parseVariable(input)
	case isAssignment(input, token):
		return parseAssignment(input, token)
	case isBinaryExpression(input, token):
		return parseBinaryExpression(input, token)
	case isCondition(input, token):
		return parseCondition(input, token)
	case isCaller(input, token):
		return parseCaller(input, token)
	}
}

func program(input lexis.TokenStream) bool {
	var prog = tokenProgram{
		Class: "program",
	}

	for !input.EOF() {
		token, err := expression(input)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			break
		}

		prog.Expression = append(prog.Expression, token)
	}

	fmt.Printf("%# v", pretty.Formatter(prog))
	return true
}
