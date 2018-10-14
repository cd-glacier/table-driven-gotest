package tdt

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	nodeparser "github.com/g-hyoga/go-node-parser"
	"github.com/google/go-cmp/cmp"
)

var ignoreInt = cmp.Comparer(func(x, y int) bool {
	return true
})

func TestFindFunc(t *testing.T) {
	tests := []struct {
		code string
		decl string
	}{
		{
			`
package main

func TestMain(t *testing.T) {}
			`,
			`
func TestMain(t *testing.T) {}
			`,
		},
		{
			`
package main

func TestMain(t *testing.T) {}
			`,
			`
func TestMain(t *testing.T) {}
			`,
		},
	}

	for _, tt := range tests {
		f, err := parser.ParseFile(token.NewFileSet(), "main.go", tt.code, parser.AllErrors)
		if err != nil {
			t.Fatalf("Failed to parseFile: %s", err.Error())
		}

		tdt := &TDT{File: f, FnName: "TestMain"}
		decl := tdt.FindFunc()

		d, err := nodeparser.ParseDecl(tt.decl)
		if err != nil {
			t.Fatalf("Failed to parseExpr: %s", err.Error())
		}
		expectedDecl, ok := (*d).(*ast.FuncDecl)
		if !ok {
			t.Fatalf("Failed to convert to *ast.FuncDecl")
		}

		if decl.Name.Name != expectedDecl.Name.Name {
			t.Fatalf("Failed to TestFindFunc. invalid decl name. actual: %s expected: %s", decl.Name.Name, expectedDecl.Name.Name)
		}

		if len(decl.Body.List) != len(expectedDecl.Body.List) {
			t.Fatalf("Failed to TestFindFunc. invalid decl body length. actual: %d expected: %d", len(decl.Body.List), len(expectedDecl.Body.List))
		}
	}
}
