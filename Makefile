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
	docker pull ubuntu:latest
	go test -v ./pkg/...

integration-test:
	ginkgo -v -r -p --skip="docker container|retraced.events" tests/ginkgo

integration-test-docker:
	docker pull ubuntu:latest
	ginkgo -v -r -p --focus="docker container" tests/ginkgo

# this task assumes a working retraced installation, and requires the following params to be set:
#
#  RETRACED_API_ENDPOINT
#  RETRACED_PROJECT_ID
#  RETRACED_API_KEY
#
# Can also optionally set
#
#  RETRACED_INSECURE_SKIP_VERIFY=1
#
integration-test-retraced:
	ginkgo -v -r -p --focus="retraced.events" tests/ginkgo

build:
	mkdir -p bin
	go build -o ./bin/support-bundle .

githooks:
	echo 'make integration-test' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	echo 'go fmt ./...' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

all: build test integration-test integration-test-docker
