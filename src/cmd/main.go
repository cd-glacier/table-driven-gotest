package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/g-hyoga/table-driven-gotest/src/logger"
	tdt "github.com/g-hyoga/table-driven-gotest/src/tabledriventest"
)

var log = logger.New()

type Option struct {
	FileName string
	FnName   string
	Index    int
}

func parseOption() *Option {
	testFileName := flag.String("f", "./src/cmd/main_test.go", "test file name")
	testFnName := flag.String("v", "TestMain", "test function name")
	testCaseIndex := flag.Int("i", 0, "test case index: start 0")
	flag.Parse()

	op := &Option{FileName: *testFileName, FnName: *testFnName, Index: *testCaseIndex}
	log.Debugf("[option] %+v", op)
	return op
}

func main() {
	op := parseOption()
	packageName, fileName := filepath.Split(op.FileName)
	packageName = packageName

	tdt, err := tdt.New(packageName, fileName, op.FnName, op.Index)
	if err != nil {
		panic(fmt.Sprintf("Failed to Tester New: %s", err.Error()))
	}

	tdt.Test()
}
