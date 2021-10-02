package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goprintffuncname",
	Doc:      "Checks that printf-like functions are named with `f` at the end.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{ // 関数定義のノードのみを訪れるようにフィルターの設定
		(*ast.FuncDecl)(nil),
	}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)

		params := funcDecl.Type.Params.List
		if len(params) < 2 {
			return
		}

		formatParamType, ok := params[len(params)-2].Type.(*ast.Ident)
		if !ok {
			return
		}

		if formatParamType.Name != "string" {
			return
		}

		argsParamType, ok := params[len(params)-1].Type.(*ast.Ellipsis)
		if !ok {
			return
		}

		elementType, ok := argsParamType.Elt.(*ast.InterfaceType)
		if !ok {
			return
		}

		if elementType.Methods != nil && len(elementType.Methods.List) != 0 {
			return
		}

		if strings.HasSuffix(funcDecl.Name.Name, "f") {
			return
		}

		pass.Reportf(node.Pos(), "printf-like formatting function '%s' should be named '%sf'",
			funcDecl.Name.Name, funcDecl.Name.Name)
	})
	return nil, nil
}
