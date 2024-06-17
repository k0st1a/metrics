// Package osexit for check direct call os.Exit in main function of main package.
package osexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check direct call os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}

		if !strings.HasSuffix(pass.Fset.Position(file.Pos()).Filename, ".go") {
			continue
		}

		ast.Inspect(file, inspect(pass))
	}

	//nolint:nilnil // default return value for analysis
	return nil, nil
}

func inspect(pass *analysis.Pass) func(ast.Node) bool {
	var isMainFunc = false

	return func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.FuncDecl:
			if v.Name.String() == "main" {
				isMainFunc = true
			} else {
				isMainFunc = false
			}
		case *ast.ExprStmt:
			isOsExitCallExp(v, pass, &isMainFunc)
		}

		return true
	}
}

func isOsExitCallExp(es *ast.ExprStmt, pass *analysis.Pass, isMainFunc *bool) {
	ce, ok := es.X.(*ast.CallExpr)
	if !ok {
		return
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	i, ok := se.X.(*ast.Ident)
	if !ok {
		return
	}

	if i.String() != "os" {
		return
	}

	if se.Sel.String() != "Exit" {
		return
	}

	if *isMainFunc {
		pass.Reportf(i.NamePos, "direct call os.Exit in main func")
	}
}
