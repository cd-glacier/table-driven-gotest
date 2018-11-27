package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/g-hyoga/table-driven-gotest/src/logger"
	tdt "github.com/g-hyoga/table-driven-gotest/src/tabledriventest"
)

var log = logger.New()

// Option represenr command line args
type Option struct {
	FileName string
	FnName   string
	Index    int
	TestCase string
}

func parseOption() *Option {
	testFileName := flag.String("file", "./src/cmd/main_test.go", "test file name")
	testFnName := flag.String("func", "TestMain", "test function name")
	testCaseIndex := flag.Int("index", 0, "test case index: start 0")
	userTestCase := flag.String("testcase", "{\"input\", \"output\"}", "test case you want to test.")
	flag.Parse()

	op := &Option{FileName: *testFileName, FnName: *testFnName, Index: *testCaseIndex, TestCase: *userTestCase}
	log.Debugf("[option] %+v", op)
	return op
}

func main() {
	op := parseOption()
	packageName, fileName := filepath.Split(op.FileName)

	tdt, err := tdt.New(packageName, fileName, op.FnName, op.Index)
	if err != nil {
		panic(fmt.Sprintf("Failed to Tester New: %s", err.Error()))
	}

	tdt.Test()
}
