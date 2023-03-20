// Package osexit определяет Analyzer, который проверяет
// наличие os.Exit() в функции main пакета main
package osexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer - анализатор, который проверяет наличие os.Exit() в функции main пакета main
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for using os.Exit() in func main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") || !strings.HasSuffix(filename, ".go") {
			continue
		}

		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.DeferStmt:
				return false
			case *ast.SelectorExpr:
				if call, ok := x.X.(*ast.Ident); ok {
					if call.Name == "os" && x.Sel.Name == "Exit" {
						pass.Reportf(call.NamePos, "using os.Exit() in func main")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
