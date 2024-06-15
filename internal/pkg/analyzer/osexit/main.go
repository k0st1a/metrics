package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for call os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, inspect(pass))
	}

	return nil, nil
}

func inspect(pass *analysis.Pass) func(ast.Node) bool {
	return func(node ast.Node) bool {
		var isMainFunc = false

		switch v := node.(type) {
		case *ast.FuncDecl:
			if v.Name.Name == "main" {
				isMainFunc = true
			} else {
				isMainFunc = false
			}
		case *ast.ExprStmt:
			isOsExitCallExp(v, pass, isMainFunc)
		}

		return true
	}
}

func isOsExitCallExp(es *ast.ExprStmt, pass *analysis.Pass, isMainFunc bool) {
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

	if i.Name != "os" {
		return
	}

	if se.Sel.Name != "Exit" {
		return
	}

	if isMainFunc {
		pass.Reportf(i.NamePos, "direct call os.Exit")
	}
}
