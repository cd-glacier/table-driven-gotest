.PHONY: help
.DEFAULT_GOAL := help

deps: ## create vendor directory with go dep.
	cd src && dep ensure

build: deps ## build main.go.
	go build -o tdt src/cmd/main.go

run: ## run main.go.
	go run ./src/cmd/main.go -f ./src/tabledriventest/tabledriventest_test.go -v TestFindFunc -i 0

test: ## unit test with gotest. ref: github.com/rakyll/gotest
	gotest -v ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
