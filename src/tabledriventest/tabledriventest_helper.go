package tabledriventest

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func getHashedDir(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	dir := strings.Replace(hash, "=", "", -1)
	return "./" + dir + "/"
}

func parseFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
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

func ReCreate(filename string) (*os.File, error) {
	err := os.Remove(filename)
	if err != nil {
		return nil, err
	}

	return os.Create(filename)
}

func splitResult(output string) string {
	testStrs := strings.Split(output, "===")

	return testStrs[1]
}

func getResult(testStr string) bool {
	return strings.Contains(strings.Split(testStr, "---")[1], "PASS")
}
