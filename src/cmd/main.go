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
	testFileName := flag.String("file", "", "test file name. example: --file ./src/cmd/main_test.go")
	testFnName := flag.String("func", "", "test function name. example: --func TestMain")
	testCaseIndex := flag.Int("index", -1, "test case index: start 0. exmaple: --index 1")
	userTestCase := flag.String("testcase", "", "test case you want to test. example: --testcase {\"input\", \"output\"}")
	flag.Parse()

	op := &Option{FileName: *testFileName, FnName: *testFnName, Index: *testCaseIndex, TestCase: *userTestCase}
	log.Debugf("[option] %+v", op)
	return op
}

func main() {
	op := parseOption()
	packageName, fileName := filepath.Split(op.FileName)

	tdt, err := tdt.New(packageName, fileName, op.FnName, op.TestCase, op.Index)
	if err != nil {
		panic(fmt.Sprintf("Failed to Tester New: %s", err.Error()))
	}

	err = tdt.Test()
	if err != nil {
		log.Errorf("Failed to Test: %s", err.Error())
		panic(err)
	}
}
