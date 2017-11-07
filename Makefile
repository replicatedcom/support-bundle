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
	go test -v ./pkg/...

integration-test:
	ginkgo -v -r -p --skip="docker container" tests/ginkgo

integration-test-docker:
	docker pull ubuntu:latest
	ginkgo -v -r -p --focus="docker container" tests/ginkgo

build:
	mkdir -p bin
	go build -o ./bin/support-bundle .

githooks:
	echo 'make integration-test' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	echo 'go fmt ./...' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

all: build test integration-test integration-test-docker
