package tdt

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type TDT struct {
	FileName string
	File     *ast.File
	FnName   string
}

func New(fileName, fnName string) (*TDT, error) {
	tdt := &TDT{}
	tdt.FileName = fileName
	tdt.FnName = fnName

	f, err := parseFile(fileName)
	if err != nil {
		return nil, err
	}
	tdt.File = f

	return tdt, nil
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
