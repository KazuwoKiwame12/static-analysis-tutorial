package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
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

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, v.fset, node)
	fmt.Printf("%s | %#v\n", buf.String(), node)

	return v
}
