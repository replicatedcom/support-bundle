SHELL := /bin/bash
.PHONY: clean deps install run test build shell all

clean:
	rm -f ./bin/support-bundle

deps:
	go install

install:
	go install

generate:
	./bin/support-bundle generate

test:
	go test -v `go list ./... | grep -v /vendor/`

integration-test:
	ginkgo -v -r -p tests/ginkgo

build:
	mkdir -p bin
	go build -o ./bin/support-bundle .

githooks:
	echo 'make integration-test' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	echo 'go fmt ./...' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

all: build test
