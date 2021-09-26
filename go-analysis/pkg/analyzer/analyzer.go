package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "goprintffuncname",
	Doc:  "Checks that printf-like functions are named with `f` at the end.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		params := funcDecl.Type.Params.List
		if len(params) != 2 {
			return true
		}
		firstParamType, ok := params[0].Type.(*ast.Ident)
		if !ok {
			return true
		}

		if firstParamType.Name != "string" {
			return true
		}

		secondParamType, ok := params[1].Type.(*ast.Ellipsis)
		if !ok {
			return true
		}

		elementType, ok := secondParamType.Elt.(*ast.InterfaceType)
		if !ok {
			return true
		}

		if elementType.Methods != nil && len(elementType.Methods.List) != 0 {
			return true
		}

		if strings.HasSuffix(funcDecl.Name.Name, "f") {
			return true
		}

		pass.Reportf(node.Pos(), "printf-like formatting function '%s' should be named '%sf'",
			funcDecl.Name.Name, funcDecl.Name.Name)
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}

	return nil, nil
}
