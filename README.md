
# go test command for table driven

You can choice a test case of table driven test and go test.
This program search test table automatically and Run go Test only test case that you select or you write in command line.

## Build 

```sh
make build
```

setting PATH in your terminal

## Usage

### Choice test case

Run test with test case you select. 

```sh
tdt -file ./src/cmd/main_test.go -func TestMain -index 2
```

#### ./src/cmd/main_test.go
```go
package main

function TestMain(t *testing.T) {
  tests := []struct {
    input string
    output string
  } {
    {"input0", "output0"},
    {"input1", "output1"},
    {"input2", "output2"}, // this test case is runed
    {"input3", "output3"},
    {"input4", "output4"},
  }

  for _, tt := range tests {
    ...
  }
}
```

### Write test case in command line

```sh
tdt -file ./src/cmd/main_test.go -func TestMain -testcase '{"input0", "output2"}'
```


