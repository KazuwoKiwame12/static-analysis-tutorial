package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "_gopher.go", nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(nil, f)
	ast.Inspect(f, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name != "Gopher" {
			return true
		}

		fmt.Println(fset.Position(ident.Pos()))
		return true
	})
}
