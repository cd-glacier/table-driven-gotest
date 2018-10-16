package tabledriventest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/g-hyoga/table-driven-gotest/src/logger"
)

type TDT struct {
	FileName      string
	File          *ast.File
	FnName        string
	TestCaseIndex int
}

var log = logger.New()

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
	fn, err := t.FindFunc()
	if err != nil {
		log.Errorf("Not Found '%s' func", t.FnName)
		panic(fmt.Sprintf("Not Found '%s' func", t.FnName))
	}
	table, err := t.FindTable(fn)
	if err != nil {
		log.Errorf("Not Found '%s' func's table", t.FnName)
		panic(fmt.Sprintf("Not Found '%s' func's table", t.FnName))
	}
	t.DeleteOtherTestCase(table)
}

func parseFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
}

func (t *TDT) FindFunc() (*ast.FuncDecl, error) {
	for _, decl := range t.File.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if "Test"+d.Name.Name == t.FnName || d.Name.Name == t.FnName {
				log.Debugf("Found '%s' func", t.FnName)
				return d, nil
			}
		}
	}
	return nil, fmt.Errorf("Not Found '%s'", t.FnName)
}

func isTable(stmt *ast.AssignStmt) bool {
	varLenIsOne := len(stmt.Rhs) == 1
	liter, isCompositeLit := stmt.Rhs[0].(*ast.CompositeLit)
	arrayType, hasArrayType := liter.Type.(*ast.ArrayType)
	_, hasStructType := arrayType.Elt.(*ast.StructType)

	log.Debugf("isTable: %t", varLenIsOne && isCompositeLit && hasArrayType && hasStructType)
	return varLenIsOne && isCompositeLit && hasArrayType && hasStructType
}

func (t *TDT) FindTable(fn *ast.FuncDecl) (*ast.AssignStmt, error) {
	for _, stmt := range fn.Body.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			if isTable(s) {
				log.Infof("Found '%s' function's table: '%s'", t.FnName, s.Lhs[0].(*ast.Ident).Name)
				return s, nil
			}
		}
	}
	return nil, fmt.Errorf("Not Found '%s' function's table", t.FnName)
}

func (t *TDT) DeleteOtherTestCase(table *ast.AssignStmt) (*ast.AssignStmt, error) {
	testCases := table.Rhs[0].(*ast.CompositeLit).Elts
	if len(testCases)-1 < t.TestCaseIndex {
		log.Errorf("Not Exist %dth index in '%s' function.", t.TestCaseIndex, t.FnName)
		return nil, fmt.Errorf("Not Exist %dth index in '%s' function", t.TestCaseIndex, t.FnName)
	}

	table.Rhs[0].(*ast.CompositeLit).Elts = testCases[t.TestCaseIndex : t.TestCaseIndex+1]
	log.Infof("Found Test Case: ", ast.Print(token.NewFileSet(), testCases[t.TestCaseIndex:t.TestCaseIndex+1]))
	return table, nil
}
