package tabledriventest

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"os/exec"

	"github.com/g-hyoga/table-driven-gotest/src/copier"
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
	tmpDirName := getHashedDir(packageName)
	copier.CopyDir(packageName, tmpDirName)

	tdt := &TDT{}
	tdt.PackageName = tmpDirName
	tdt.FileName = tdt.PackageName + fileName
	tdt.FnName = fnName

	f, err := parseFile(tdt.FileName)
	if err != nil {
		log.Errorf("Failed to tdt New. file name: %s", tdt.FileName, err.Error())
		defer os.RemoveAll(tmpDirName)
		return nil, err
	}
	tdt.File = f
	tdt.TestCaseIndex = index

	log.Debugf("Package: %s, file: %s, function: %s, index: %d", tdt.PackageName, tdt.FileName, tdt.FnName, tdt.TestCaseIndex)
	return tdt, nil
}

func (t *TDT) Test() error {
	defer os.RemoveAll(t.PackageName)
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

func test(packageName, fnName string) bool {
	cmd := exec.Command("go", "test", "-v", packageName, "-run", fnName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("%s", err.Error())
	}

	fmt.Println(string(out))
	return getResult(splitResult(string(out)))
}

func isTable(stmt *ast.AssignStmt) bool {
	varLenIsOne := len(stmt.Rhs) == 1
	liter, isCompositeLit := stmt.Rhs[0].(*ast.CompositeLit)
	arrayType, hasArrayType := liter.Type.(*ast.ArrayType)
	_, hasStructType := arrayType.Elt.(*ast.StructType)

	log.Debugf("isTable: %t", varLenIsOne && isCompositeLit && hasArrayType && hasStructType)
	return varLenIsOne && isCompositeLit && hasArrayType && hasStructType
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
