package tabledriventest

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

func TestFindTable(t *testing.T) {
	tests := []struct {
		fnCode string
		table  string
	}{
		{
			`
func TestMain(t *testing.T) {
	tests := []struct{
		input string
		output string
	}{
		{
			"input0",
			"output0",
		},
		{
			"input1",
			"output1",
		},
		{
			"input2",
			"output2",
		},
	}
}
			`,
			`
tests := []struct{
	input string
	output string
}{
	{
		"input0",
		"output0",
	},
	{
		"input1",
		"output1",
	},
	{
		"input2",
		"output2",
	},
}
			`,
		},
	}

	for _, tt := range tests {
		tdt := &TDT{}

		d, err := nodeparser.ParseDecl(tt.fnCode)
		if err != nil {
			t.Fatalf("Failed to ParseDecl: %s", err.Error())
		}
		decl, ok := (*d).(*ast.FuncDecl)
		if !ok {
			t.Fatalf("Failed to convert to *ast.FuncDecl")
		}
		table := tdt.FindTable(decl)

		stmt, err := nodeparser.ParseStmt(tt.table)
		if err != nil {
			t.Fatalf("Failed to ParseExpr: %s", err.Error())
		}

		expectedTable, ok := (*stmt).(*ast.AssignStmt)
		if !ok {
			t.Fatalf("Failed to convert to *ast.Stmt")
		}

		compositeLit, ok := table.Rhs[0].(*ast.CompositeLit)
		if !ok {
			t.Fatalf("Failed to convert to *ast.CompositeLit")
		}

		expectedCompositeLit, ok := expectedTable.Rhs[0].(*ast.CompositeLit)
		if !ok {
			t.Fatalf("Failed to convert to *ast.CompositeLit")
		}

		if len(compositeLit.Elts) != len(expectedCompositeLit.Elts) {
			t.Fatalf("Invalid table elements length. actual: %d, expected: %d", len(compositeLit.Elts), len(expectedCompositeLit.Elts))
		}
	}
}

func TestDeleteOtherTestCase(t *testing.T) {
	tests := []struct {
		table        string
		index        int
		tableOnlyOne string
	}{
		{
			`
tests := []struct{
	input string
	output string
}{
	{
		"input0",
		"output0",
	},
	{
		"input1",
		"output1",
	},
	{
		"input2",
		"output2",
	},
}
			`,
			1,
			`
tests := []struct{
	input string
	output string
}{
	{
		"input1",
		"output1",
	},
}
			`,
		},
	}

	for _, tt := range tests {
		stmt, err := nodeparser.ParseStmt(tt.table)
		if err != nil {
			t.Fatalf("Failed to ParseExpr: %s", err.Error())
		}
		table, ok := (*stmt).(*ast.AssignStmt)
		if !ok {
			t.Fatalf("Failed to convert to *ast.AssignStmt")
		}

		tdt := &TDT{TestCaseIndex: tt.index}

		newTable := tdt.DeleteOtherTestCase(table)

		if !isTable(newTable) {
			t.Fatalf("Failed to DeleteOtherTestCase: table definition is not full")
		}

		expectedElts := table.Rhs[0].(*ast.CompositeLit).Elts
		elts := newTable.Rhs[0].(*ast.CompositeLit).Elts

		if len(expectedElts) != len(elts) {
			t.Fatalf("Failed to DeleteOtherTestCase. test case length is not valid. actual: %d, expected: %d", len(elts), len(expectedElts))
		}

		for i, e := range elts {
			liter, ok := e.(*ast.CompositeLit)
			if !ok {
				t.Fatalf("Failed to convert to *ast.CompositeLit")
			}
			expectedLiter, ok := expectedElts[i].(*ast.CompositeLit)
			if !ok {
				t.Fatalf("Failed to convert to *ast.CompositeLit")
			}

			key, ok := liter.Elts[0].(*ast.BasicLit)
			if !ok {
				t.Fatalf("Failed to convert to *ast.BasicLit")
			}
			expectedKey, ok := expectedLiter.Elts[0].(*ast.BasicLit)
			if !ok {
				t.Fatalf("Failed to convert to *ast.BasicLit")
			}

			if key.Value != expectedKey.Value {
				t.Fatalf("Failed to DeleteOtherTestCase. invalid key. actual: %s, expected: %s", key.Value, expectedKey.Value)
			}

			value, ok := liter.Elts[1].(*ast.BasicLit)
			if !ok {
				t.Fatalf("Failed to convert to *ast.BasicLit")
			}
			expectedValue, ok := expectedLiter.Elts[1].(*ast.BasicLit)
			if !ok {
				t.Fatalf("Failed to convert to *ast.BasicLit")
			}

			if value.Value != expectedValue.Value {
				t.Fatalf("Failed to DeleteOtherTestCase. invalid value. actual: %s, expected: %s", value.Value, expectedValue.Value)
			}

		}
	}
}
