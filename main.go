package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	for _, filepath := range os.Args[1:] {
		if filepath == "--" {
			continue
		}

		f, err := parser.ParseFile(v.fset, filepath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filepath, err)
		}

		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

// example.goのmyLogのlintする関数
// 関数の末尾が"f"で終了していない関数を見つける
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	// 関数の定義を取得できるか
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	// 関数の引数が2個か
	params := funcDecl.Type.Params.List
	if len(params) != 2 {
		return v
	}

	// 関数の1番目の引数が識別子か = 変数か確認している(= _のblankでないことの確認)
	// https://golang.org/ref/spec#Identifiers
	firstParamType, ok := params[0].Type.(*ast.Ident)
	if !ok {
		return v
	}

	// 関数の1番目の引数の型がstringか
	if firstParamType.Name != "string" {
		return v
	}

	// 関数の2番目の引数が"..."の形式か
	secondParamType, ok := params[1].Type.(*ast.Ellipsis)
	if !ok {
		return v
	}

	// 関数の2番目の引数の型がinterface型か
	elementType, ok := secondParamType.Elt.(*ast.InterfaceType)
	if !ok {
		return v
	}

	// 関数の2番目の引数のinterface型が、empty interfaceか
	if elementType.Methods != nil && len(elementType.Methods.List) != 0 {
		return v
	}

	// 関数の末尾に"f"がついているか
	if strings.HasSuffix(funcDecl.Name.Name, "f") {
		return v
	}

	fmt.Printf("%s: printf-like formatting function '%s' should be named '%sf'\n",
		v.fset.Position(node.Pos()), funcDecl.Name.Name, funcDecl.Name.Name,
	)

	return v
}
