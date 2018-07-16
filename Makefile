.PHONY: docker deps fmt vet _vet lint _lint test _test build _build build-deps dep-deps citest ci-upload-coverage e2e e2e-supportbundle e2e-supportbundle-core e2e-supportbundle-docker e2e-supportbundle-swarm ci-e2e ci-e2e-supportbundle ci-e2e-supportbundle-core ci-e2e-supportbundle-docker ci-e2e-supportbundle-swarm goreleaser

SHELL := /bin/bash
SRC = $(shell find . -name "*.go")
PKG := github.com/replicatedhq/support-bundle
VERSION := $(shell git describe --tags --always --dirty)
SHA := $(shell git log --pretty=format:'%H' -n 1)
ARCH ?= amd64
ifeq ($(shell uname), Darwin)
	BUILD_TIME := $(shell date -u +%FT%T)
else
	BUILD_TIME := $(shell date --rfc-3339=seconds | sed 's/ /T/')
endif
#paths within WSL start with /mnt/c/...
#docker does not recognize this fact
#this strips the first 5 characters (leaving /c/...) if the kernel releaser is Microsoft
ifeq ($(shell uname -r | tail -c 10), Microsoft)
	BUILD_DIR := $(shell pwd | cut -c 5-)
else
	BUILD_DIR := $(shell pwd)
endif

docker:
	docker build -t support-bundle .

deps:
	dep ensure -v; dep prune -v

fmt:
	goimports -w pkg
	goimports -w cmd

vet: fmt _vet

_vet:
	go vet ./pkg/...
	go vet ./cmd/...

lint: vet _lint

_lint:
	golint ./pkg/... | grep -v "should have comment" | grep -v "comment on exported" || :
	golint ./cmd/... | grep -v "should have comment" | grep -v "comment on exported" || :

test: lint _test

_test:
	go test ./pkg/...

build: test _build

_build: bin/support-bundle

bin/support-bundle: $(SRC)
	go build \
		-i \
		-o bin/support-bundle \
		-ldflags " \
		-X $(PKG)/pkg/version.version=$(VERSION) \
		-X $(PKG)/pkg/version.gitSHA=$(SHA) \
		-X $(PKG)/pkg/version.buildTime=$(BUILD_TIME) \
		" \
		./cmd/support-bundle
	@echo built bin/support-bundle

build-deps:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/jteeuwen/go-bindata/go-bindata

dep-deps:
	go get github.com/golang/dep/cmd/dep

.state/coverage.out: $(SRC)
	@mkdir -p .state/
	go test -coverprofile=.state/coverage.out -v ./pkg/...

citest: lint .state/coverage.out

.state/cc-test-reporter:
	@mkdir -p .state/
	wget -O .state/cc-test-reporter https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64
	chmod +x .state/cc-test-reporter

ci-upload-coverage: .state/coverage.out .state/cc-test-reporter
	./.state/cc-test-reporter format-coverage -o .state/codeclimate/codeclimate.json -t gocov .state/coverage.out
	./.state/cc-test-reporter upload-coverage -i .state/codeclimate/codeclimate.json

e2e: e2e-supportbundle

e2e-supportbundle: e2e-supportbundle-core e2e-supportbundle-docker

e2e-supportbundle-core:
	@docker run                                                             \
		-ti                                                                 \
		--rm                                                                \
		-v "$(BUILD_DIR):/go/src/$(PKG)"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
		-w /go/src/$(PKG)                                                   \
		-l com.replicated.support-bundle=true                               \
		golang:1.10                                                         \
		/bin/sh -c "                                                        \
			./e2e/collect/e2e.sh                                            \
		"

e2e-supportbundle-docker:
	docker pull ubuntu:16.04
	@docker run                                                             \
		-ti                                                                 \
		--rm                                                                \
		-v "$(BUILD_DIR):/go/src/$(PKG)"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
		-w /go/src/$(PKG)                                                   \
		-l com.replicated.support-bundle=true                               \
		-e DOCKER=1                                                         \
		golang:1.10                                                         \
		/bin/sh -c "                                                        \
			./e2e/collect/e2e.sh                                            \
		"

e2e-supportbundle-swarm:
	@docker run                                                             \
		-ti                                                                 \
		--rm                                                                \
		-v "$(BUILD_DIR):/go/src/$(PKG)"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
		-w /go/src/$(PKG)                                                   \
		-l com.replicated.support-bundle=true                               \
		-e SWARM=1                                                          \
		golang:1.10                                                         \
		/bin/sh -c "                                                        \
			./e2e/collect/e2e.sh                                            \
		"

ci-e2e: ci-e2e-supportbundle

ci-e2e-supportbundle: ci-e2e-supportbundle-core ci-e2e-supportbundle-docker

ci-e2e-supportbundle-core:
	./e2e/collect/e2e.sh

ci-e2e-supportbundle-docker:
	docker pull ubuntu:16.04
	DOCKER=true ./e2e/collect/e2e.sh

ci-e2e-supportbundle-swarm:
	SWARM=true ./e2e/collect/e2e.sh

goreleaser: .state/goreleaser

.state/goreleaser: .goreleaser.unstable.yml deploy/Dockerfile-collect $(SRC)
	@mkdir -p .state
	@touch .state/goreleaser
	curl -sL https://git.io/goreleaser | bash -s -- --snapshot --rm-dist --config .goreleaser.unstable.yml
