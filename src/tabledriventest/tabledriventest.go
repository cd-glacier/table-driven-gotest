package tabledriventest

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"github.com/g-hyoga/table-driven-gotest/src/logger"
)

type TDT struct {
	PackageName   string
	FileName      string
	File          *ast.File
	FnName        string
	TestCaseIndex int
	Passed        bool
}

var log = logger.New()

func New(packageName, fileName, fnName string, index int) (*TDT, error) {
	tdt := &TDT{}
	tdt.PackageName = "./" + packageName
	tdt.FileName = tdt.PackageName + fileName
	tdt.FnName = fnName

	f, err := parseFile(tdt.FileName)
	if err != nil {
		log.Errorf("Failed to tdt New. file name: %s", tdt.FileName, err.Error())
		return nil, err
	}
	tdt.File = f
	tdt.TestCaseIndex = index

	log.Debugf("Package: %s, file: %s, function: %s, index: %d", tdt.PackageName, tdt.FileName, tdt.FnName, tdt.TestCaseIndex)
	return tdt, nil
}

func ReCreate(filename string) (*os.File, error) {
	err := os.Remove(filename)
	if err != nil {
		return nil, err
	}

	return os.Create(filename)
}

func (t *TDT) Test() error {
	fn, err := t.FindFunc()
	if err != nil {
		log.Errorf("Not Found '%s' func", t.FnName)
		return err
	}

	table, err := t.FindTable(fn)
	if err != nil {
		log.Errorf("Not Found '%s' func's table", t.FnName)
		return err
	}

	_, err = t.DeleteOtherTestCase(table)
	if err != nil {
		log.Errorf("Not Found %dth test case", t.TestCaseIndex)
		return err
	}

	file, err := ReCreate(t.FileName)
	defer file.Close()
	if err != nil {
		log.Error("Failed to open %s: %s", t.FileName, err.Error())
		return err
	}
	err = format.Node(file, token.NewFileSet(), t.File)
	if err != nil {
		log.Errorf("Failed to format.Node: %s", err.Error())
		return err
	}

	t.Passed = test(t.PackageName, t.FnName)

	log.Debugf("Package: %s, file: %s, function: %s, index: %d, passed: %t", t.PackageName, t.FileName, t.FnName, t.TestCaseIndex, t.Passed)
	return nil
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
	log.Infof("Found %dth Test Case", t.TestCaseIndex)
	// log.Debugf("Found Test Case: ", ast.Print(token.NewFileSet(), testCases[t.TestCaseIndex:t.TestCaseIndex+1]))
	return table, nil
}

func splitResult(output string) string {
	testStrs := strings.Split(output, "===")

	return testStrs[1]
}

func getResult(testStr string) bool {
	return strings.Contains(strings.Split(testStr, "---")[1], "PASS")
}

func test(packageName, fnName string) bool {
	cmd := exec.Command("go", "test", "-v", packageName, "-run", fnName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("%s", err.Error())
	}

	fmt.Println(string(out))
	return getResult(splitResult(string(out)))
}
