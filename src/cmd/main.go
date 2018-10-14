package main

import (
	"flag"

	"github.com/k0kubun/pp"
)

func main() {
	packageName := *flag.String("package", "./src/cmd/", "package name")
	testFileName := *flag.String("file", "main_test.go", "test file name")
	testFnName := *flag.String("function", "TestMain", "test function name")
	testCaseNum := *flag.Int("testcase", 0, "test case number: start 0")
	flag.Parse()

	pp.Println(packageName, testFileName, testFnName, testCaseNum)
}
