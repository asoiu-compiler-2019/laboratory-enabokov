package semantics

import (
	"fmt"
	"log"
	"strings"

	"github.com/enabokov/language/syntax"
)

type mapVariables map[string]string
type mapFunctions map[string][]string
type mapImports map[string]int
type mapPackage map[string]int

func scanCall(call syntax.TokenCall, localVars []string, funcName string, funcs *mapFunctions, imports *mapImports) error {
	res := strings.Split(call.Func.Name, `.`)
	mainImport, fn := res[0], res[1]

	// check if function defined
	if mainImport == `locals` {
		if _, ok := (*funcs)[fn]; !ok {
			return fmt.Errorf("No such method %s", fn)
		}

		if len(call.Args) != len((*funcs)[fn]) {
			return fmt.Errorf("Wrong call %s. Too many arguments: expect %v - got %v", call.Func.Name, (*funcs)[fn], call.Args)
		}
	} else {
		// check if there was import for this function
		if _, ok := (*imports)[mainImport]; !ok {
			return fmt.Errorf("No such import %s. Failed on %s", mainImport, call.Func.Name)
		}
	}

	// check if passed args are defined earlier
	var matched []string
	args := (*funcs)[funcName]
	for _, v := range localVars {
		for _, arg := range args {
			if v == arg {
				matched = append(matched, v)
			}
		}
	}

	if len(args) != len(matched) {
		return fmt.Errorf("Wrong call %s", call.Func.Name)
	}

	return nil
}

func scanVariable(variable syntax.TokenVariable, localVars *[]string, funcName string) error {
	for _, v := range *localVars {
		if variable.Name == v {
			return fmt.Errorf("Variable `%s` is already defined in `%s`", variable.Name, funcName)
		}
	}

	*localVars = append(*localVars, variable.Name)
	return nil
}

func scanBody(body []syntax.ASTNode, localVars []string, funcName string, funcs *mapFunctions, imports *mapImports) (err error) {
	for _, expr := range body {
		astCall, ok := expr.(syntax.TokenCall)
		if ok {
			err = scanCall(astCall, localVars, funcName, funcs, imports)
			if err != nil {
				return err
			}

			continue
		}

		astVar, ok := expr.(syntax.TokenVariable)
		if ok {
			err = scanVariable(astVar, &localVars, funcName)
			if err != nil {
				return err
			}

			continue
		}
	}

	return err
}

func scanFunction(ast syntax.TokenFunction, funcs *mapFunctions, imports *mapImports) (err error) {
	var localVariables []string
	var params []string

	for _, param := range ast.Params {
		params = append(params, param.Name)
		localVariables = append(localVariables, param.Name)
	}

	// register in global scope
	(*funcs)[ast.Name] = params

	return scanBody(ast.Body, localVariables, ast.Name, funcs, imports)
}

func scanImport(ast syntax.TokenImport, funcs *mapFunctions, imports *mapImports) error {
	if _, ok := (*imports)[ast.Value]; ok {
		return fmt.Errorf("Import already exist `%s`", ast.Value)
	}

	(*imports)[ast.Value] = 1
	return nil
}

func scanPackage(ast syntax.TokenPackage, packages *mapPackage) error {
	if _, ok := (*packages)[ast.Value]; ok {
		return fmt.Errorf("Package can be defined only once `%s`", ast.Value)
	}

	return nil
}

func scan(ast syntax.TokenProgram) (err error) {
	functions := make(mapFunctions)
	imports := make(mapImports)
	packages := make(mapPackage)

	for _, expr := range ast.Expression {
		astFunc, ok := expr.(syntax.TokenFunction)
		if ok {
			err = scanFunction(astFunc, &functions, &imports)
			if err != nil {
				log.Fatalln(err)
			}
			continue
		}

		astImport, ok := expr.(syntax.TokenImport)
		if ok {
			err = scanImport(astImport, &functions, &imports)
			if err != nil {
				return err
			}
			continue
		}

		astPackage, ok := expr.(syntax.TokenPackage)
		if ok {
			err = scanPackage(astPackage, &packages)
			if err != nil {
				return err
			}
			continue
		}
	}

	return err
}
