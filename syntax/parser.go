package syntax

import (
	"fmt"
	"log"
	"strconv"

	"github.com/enabokov/language/lexis"
)

func parseParams(input lexis.TokenStream, params *[]TokenVariable, requiredParam bool) error {
	nextToken := input.Peek()
	if nextToken.Class == lexis.ClassVariable {
		if !requiredParam {
			return input.Croak(fmt.Sprintf("Got `%s`. Expected `,` or `)`", input.Peek().Value))
		}

		t := input.Next()

		*params = append(
			*params,
			TokenVariable{
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

func parseArgs(input lexis.TokenStream, params *[]TokenVariable, requiredParam bool) error {
	nextToken := input.Peek()
	if nextToken.Class == lexis.ClassVariable || nextToken.Class == lexis.ClassString {
		if !requiredParam {
			return input.Croak(fmt.Sprintf("Got `%s`. Expected `,` or `)`", input.Peek().Value))
		}

		t := input.Next()
		*params = append(
			*params,
			TokenVariable{
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

func parseBody(input lexis.TokenStream, body *[]ASTNode) error {
	token := input.Next()
	if token.Class != lexis.ClassPunctuation || token.Value != `{` {
		return input.Croak(fmt.Sprintf("Got `%s`. Expected {", token.Value))
	}

	for input.Peek().Class != lexis.ClassPunctuation || input.Peek().Value != `}` {
		ASTNode, err := expression(input)
		if err != nil {
			return err
		}

		*body = append(*body, ASTNode)
	}

	input.Next()
	return nil
}

func parseFunction(input lexis.TokenStream) (TokenFunction, error) {
	// get function name
	name := input.Next().Value

	var (
		params []TokenVariable
		err    error
	)

	if token := input.Next(); token.Class != lexis.ClassPunctuation || token.Value != `(` {
		return TokenFunction{}, input.Croak(fmt.Sprintf("Got `%s`. Expected `(`", token.Value))
	}

	if err = parseParams(input, &params, true); err != nil {
		return TokenFunction{}, err
	}

	var body []ASTNode
	err = parseBody(input, &body)

	return TokenFunction{
		Class:  `function`,
		Name:   name,
		Params: params,
		Body:   body,
	}, err
}

func parsePackage(input lexis.TokenStream) (TokenPackage, error) {
	return TokenPackage{Class: `package`, Value: input.Next().Value}, nil
}

func parseVariable(input lexis.TokenStream) (TokenVariable, error) {
	name := input.Next().Value
	token := input.Next()
	if token.Class != lexis.ClassType {
		return TokenVariable{}, input.Croak(fmt.Sprintf("Got `%s`. Expected type of variable", token.Value))
	}

	return TokenVariable{Class: `variable`, Name: name, Type: token.Value}, nil
}

func parseImport(input lexis.TokenStream) (TokenImport, error) {
	return TokenImport{Class: `import`, Value: input.Next().Value}, nil
}

func parseCaller(input lexis.TokenStream, token *lexis.Token) (TokenCall, error) {
	nextToken := input.Next()
	if input.Peek().Class == lexis.ClassPunctuation && input.Peek().Value == "(" {
		input.Next()
	}

	var args []TokenVariable

	err := parseArgs(input, &args, true)
	return TokenCall{
		Class: `caller`,
		Func:  TokenVariable{Class: token.Class, Name: token.Value + nextToken.Value},
		Args:  args,
	}, err
}

func _parse(operators []string, numbers []string, result *TokenBinaryOrAssign) error {
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
		result.Left = TokenPrimitive{
			Class: `number`,
			Value: num,
		}
	}

	result.Class = `binary`
	result.Right = &TokenBinaryOrAssign{}
	return _parse(operators, numbers, result.Right)
}

func parseBinaryExpression(input lexis.TokenStream, token *lexis.Token) (TokenBinaryOrAssign, error) {
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

	var res = TokenBinaryOrAssign{Class: `assignment`}
	err := _parse(operators, numbers, &res)
	return res, err
}

func parseAssignment(input lexis.TokenStream, token *lexis.Token) (TokenBinaryOrAssign, error) {
	nextToken := input.Next()
	number := input.Peek()

	res, err := parseBinaryExpression(input, number)

	if _, err := strconv.Atoi(token.Value); err == nil {
		return TokenBinaryOrAssign{}, input.Croak(fmt.Sprintf("Got %s. Cannot assign to number", token.Value))
	}

	return TokenBinaryOrAssign{
		Class:    `assignment`,
		Operator: nextToken.Value,
		Left:     TokenPrimitive{Class: `variable`, Value: token.Value},
		Right:    &res,
	}, err
}

func parseCondition(input lexis.TokenStream, token *lexis.Token) (TokenCondition, error) {
	input.Next()

	var body []ASTNode
	cond, err := parseBinaryExpression(input, token)
	err = parseBody(input, &body)

	return TokenCondition{
		Class:     `condition`,
		Condition: cond,
		Do:        body,
	}, err
}

func expression(input lexis.TokenStream) (ASTNode, error) {
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

func program(input lexis.TokenStream) (ast TokenProgram) {
	for !input.EOF() {
		token, err := expression(input)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			break
		}

		ast.Expression = append(ast.Expression, token)
	}

	ast.Class = `program`
	return ast
}
