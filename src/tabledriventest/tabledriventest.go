package tabledriventest

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type TDT struct {
	FileName      string
	File          *ast.File
	FnName        string
	TestCaseIndex int
}

func New(fileName, fnName string, index int) (*TDT, error) {
	tdt := &TDT{}
	tdt.FileName = fileName
	tdt.FnName = fnName

	f, err := parseFile(fileName)
	if err != nil {
		return nil, err
	}
	tdt.File = f
	tdt.TestCaseIndex = index

	return tdt, nil
}

func (t *TDT) Test() {

}

func parseFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
}

func (t *TDT) FindFunc() *ast.FuncDecl {
	for _, decl := range t.File.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name.Name == t.FnName {
				return d
			}
		}
	}
	return nil
}

func isTable(stmt *ast.AssignStmt) bool {
	varLenIsOne := len(stmt.Rhs) == 1
	liter, isCompositeLit := stmt.Rhs[0].(*ast.CompositeLit)
	arrayType, hasArrayType := liter.Type.(*ast.ArrayType)
	_, hasStructType := arrayType.Elt.(*ast.StructType)

	return varLenIsOne && isCompositeLit && hasArrayType && hasStructType
}

func (t *TDT) FindTable(fn *ast.FuncDecl) *ast.AssignStmt {
	for _, stmt := range fn.Body.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			if isTable(s) {
				return s
			}
		}
	}
	return nil
}

func (t *TDT) DeleteOtherTestCase(table *ast.AssignStmt) *ast.AssignStmt {
	testCases := table.Rhs[0].(*ast.CompositeLit).Elts
	table.Rhs[0].(*ast.CompositeLit).Elts = testCases[t.TestCaseIndex : t.TestCaseIndex+1]
	return table
}
