SHELL := /bin/bash
.PHONY: clean deps install run test build shell all

SUPPORTBUNDLE_VERSION?=1.0.0

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
	ginkgo -v -r -p --skip="docker container|retraced.events|journald.logs" tests/ginkgo

integration-test-docker:
	docker pull ubuntu:latest
	ginkgo -v -r -p --focus="docker container|journald.logs" tests/ginkgo

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
	go build \
		-ldflags=" \
		-X github.com/replicatedcom/support-bundle/version.version=$(SUPPORTBUNDLE_VERSION) \
		-X github.com/replicatedcom/support-bundle/version.gitSHA=$(shell git log --pretty=format:'%h' -n 1) \
		-X github.com/replicatedcom/support-bundle/version.buildTime=$(shell date --rfc-3339=seconds | sed 's/ /T/')" \
		-o ./bin/support-bundle .

githooks:
	echo 'make integration-test' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	echo 'go fmt ./...' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

all: build test integration-test integration-test-docker
