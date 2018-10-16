package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/g-hyoga/table-driven-gotest/src/copier"
	"github.com/g-hyoga/table-driven-gotest/src/logger"
	tdt "github.com/g-hyoga/table-driven-gotest/src/tabledriventest"
)

var log = logger.New()

func hash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	dir := strings.Replace(hash, "=", "", -1)
	return dir
}

func main() {
	packageName := flag.String("package", "./src/cmd/", "package name")
	testFileName := flag.String("file", "main_test.go", "test file name")
	testFnName := flag.String("function", "TestMain", "test function name")
	testCaseIndex := flag.Int("index", 0, "test case index: start 0")
	flag.Parse()

	log.Debugf("[option] package: %s, file: %s, function: %s, testcase index: %d", *packageName, *testFileName, *testFnName, *testCaseIndex)

	tmpDirName := hash(*packageName) + "/"
	defer os.RemoveAll(tmpDirName)
	copier.CopyDir(*packageName, tmpDirName)

	tdt, err := tdt.New(tmpDirName, *testFileName, *testFnName, *testCaseIndex)
	if err != nil {
		panic(fmt.Sprintf("Failed to Tester New: %s", err.Error()))
	}

	tdt.Test()
}
