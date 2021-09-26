package main

import (
	"go-analysis/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
