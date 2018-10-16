package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"

	"github.com/g-hyoga/table-driven-gotest/src/copier"
	"github.com/k0kubun/pp"
)

func hash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func main() {
	packageName := *flag.String("package", "./src/cmd/", "package name")
	testFileName := *flag.String("file", "main_test.go", "test file name")
	testFnName := *flag.String("function", "TestMain", "test function name")
	testCaseNum := *flag.Int("testcase", 0, "test case number: start 0")
	flag.Parse()

	pp.Println(packageName, testFileName, testFnName, testCaseNum)

	tmpDirName := hash(packageName)
	copier.CopyDir(packageName, tmpDirName)
}
