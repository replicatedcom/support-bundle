.PHONY: docs docker fmt vet _vet lint _lint test _test build _build bindata _mockgen mockgen build-deps ci-test ci-upload-coverage e2e e2e-analyze e2e-supportbundle e2e-supportbundle-core e2e-supportbundle-docker e2e-supportbundle-swarm ci-e2e ci-e2e-supportbundle ci-e2e-supportbundle-core ci-e2e-supportbundle-docker ci-e2e-supportbundle-swarm goreleaser

SHELL := /bin/bash
SRC = $(shell find . -name "*.go")
PKG := github.com/replicatedcom/support-bundle
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
DOCKER_REPO ?= replicated

docs:
	make -C hack/docs pipeline-nointegration

docker:
	docker build -t support-bundle .

shell:
	docker run --rm -it --name support-bundle \
		-v `pwd`:/go/src/github.com/replicatedcom/support-bundle \
		support-bundle \
		bash

fmt:
	goimports -w pkg
	goimports -w cmd

vet: fmt _vet

_vet:
	go vet ./pkg/... ./cmd/...

lint: vet _lint

_lint:
	golint ./pkg/... ./cmd/... \
		| grep -v "should have comment" \
		| grep -v "comment on exported" \
		|| :

test: lint _test

_test: bindata
	go test -race ./pkg/... ./cmd/...

build: test _build

_build: bin/analyze bin/support-bundle

bindata: pkg/analyze/api/v1/defaultspec/asset.go pkg/collect/bundle/defaultspec/asset.go

pkg/analyze/api/v1/defaultspec/asset.go: pkg/analyze/api/v1/defaultspec/assets/*
	go-bindata \
		-pkg defaultspec \
		-prefix pkg/analyze/api/v1/defaultspec/ \
		-o pkg/analyze/api/v1/defaultspec/asset.go \
		pkg/analyze/api/v1/defaultspec/assets/

pkg/collect/bundle/defaultspec/asset.go: pkg/collect/bundle/defaultspec/assets/*
	go-bindata \
		-pkg defaultspec \
		-prefix pkg/collect/bundle/defaultspec/ \
		-o pkg/collect/bundle/defaultspec/asset.go \
		pkg/collect/bundle/defaultspec/assets/

_mockgen:
	rm -rf pkg/test-mocks
	mockgen \
		-destination pkg/test-mocks/collect/bundle/reader/bundlereader.go \
		-package reader \
		github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader \
		BundleReader
	mockgen \
		-destination pkg/test-mocks/collect/bundle/reader/scanner.go \
		-package reader \
		github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader \
		Scanner

mockgen: _mockgen fmt

bin/analyze: $(SRC) pkg/analyze/api/v1/defaultspec/asset.go
	CGO_ENABLED=0 go build \
		-ldflags " -s -w \
		-X $(PKG)/pkg/version.version=$(VERSION) \
		-X $(PKG)/pkg/version.gitSHA=$(SHA) \
		-X $(PKG)/pkg/version.buildTime=$(BUILD_TIME) \
		" \
		-o bin/analyze \
		./cmd/analyze
	@echo built bin/analyze

bin/support-bundle: $(SRC) pkg/collect/bundle/defaultspec/asset.go
	CGO_ENABLED=0 go build \
		-ldflags " -s -w \
		-X $(PKG)/pkg/version.version=$(VERSION) \
		-X $(PKG)/pkg/version.gitSHA=$(SHA) \
		-X $(PKG)/pkg/version.buildTime=$(BUILD_TIME) \
		" \
		-o bin/support-bundle \
		./cmd/support-bundle
	@echo built bin/support-bundle

build-deps:
	go install golang.org/x/lint/golint@latest
	go install golang.org/x/tools/cmd/goimports@v0.25.0
	go install github.com/a-urth/go-bindata/go-bindata@latest
	go install github.com/onsi/ginkgo/ginkgo@v1.16.5
	go install github.com/golang/mock/mockgen@latest

.state/coverage.out: $(SRC)
	@mkdir -p .state/
	go test -race -coverprofile=.state/coverage.out -v ./pkg/...

ci-test: lint .state/coverage.out

.state/cc-test-reporter:
	@mkdir -p .state/
	wget -O .state/cc-test-reporter https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64
	chmod +x .state/cc-test-reporter

ci-upload-coverage: .state/coverage.out .state/cc-test-reporter
	./.state/cc-test-reporter format-coverage -o .state/codeclimate/codeclimate.json -t gocov .state/coverage.out
	./.state/cc-test-reporter upload-coverage -i .state/codeclimate/codeclimate.json

e2e: e2e-analyze e2e-supportbundle

e2e-analyze:
	ginkgo -v -r -p e2e/analyze

e2e-supportbundle: e2e-supportbundle-core e2e-supportbundle-docker

e2e-supportbundle-core:
	@docker run                                                             \
		-it                                                                 \
		--rm                                                                \
		-v "$(BUILD_DIR):/go/src/$(PKG)"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
		-w /go/src/$(PKG)                                                   \
		-l com.replicated.support-bundle=true                               \
		golang:1.23                                                         \
		/bin/sh -c "                                                        \
			./e2e/collect/e2e.sh                                            \
		"

e2e-supportbundle-docker:
	docker pull ubuntu:16.04
	@docker run                                                             \
		-it                                                                 \
		--rm                                                                \
		-v "$(BUILD_DIR):/go/src/$(PKG)"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
		-w /go/src/$(PKG)                                                   \
		-l com.replicated.support-bundle=true                               \
		-e DOCKER=1                                                         \
		golang:1.23                                                         \
		/bin/sh -c "                                                        \
			./e2e/collect/e2e.sh                                            \
		"

e2e-supportbundle-swarm:
	@docker run                                                             \
		-it                                                                 \
		--rm                                                                \
		-v "$(BUILD_DIR):/go/src/$(PKG)"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
		-w /go/src/$(PKG)                                                   \
		-l com.replicated.support-bundle=true                               \
		-e SWARM=1                                                          \
		golang:1.23                                                         \
		/bin/sh -c "                                                        \
			./e2e/collect/e2e.sh                                            \
		"

ci-e2e: e2e-analyze ci-e2e-supportbundle

ci-e2e-supportbundle: ci-e2e-supportbundle-core ci-e2e-supportbundle-docker

ci-e2e-supportbundle-core:
	./e2e/collect/e2e.sh

ci-e2e-supportbundle-docker:
	docker pull ubuntu:16.04
	DOCKER=true ./e2e/collect/e2e.sh

ci-e2e-supportbundle-swarm:
	SWARM=true ./e2e/collect/e2e.sh

goreleaser: .state/goreleaser

.state/goreleaser: .state/base .goreleaser.unstable.yml deploy/Dockerfile-analyze deploy/Dockerfile-collect $(SRC)
	@mkdir -p .state
	curl -sL https://git.io/goreleaser | bash -s -- --snapshot --rm-dist --config .goreleaser.unstable.yml
	@touch .state/goreleaser

support-bundle-generate: goreleaser
	@docker run \
		-it \
		--rm \
		--name support-bundle \
		--volume $(BUILD_DIR):/out \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--env LOG_LEVEL=DEBUG \
		--pid host \
		--workdir /out  \
		$(DOCKER_REPO)/support-bundle:alpha \
		generate

.state/base: deploy/Dockerfile-base
	@mkdir -p .state
	docker build --pull -t replicated/support-bundle:base -f deploy/Dockerfile-base deploy/
	@touch .state/base

grype:
	curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b .

build-base:
	docker build --pull -t replicated/support-bundle:base -f deploy/Dockerfile-base deploy/

scan-base: build-base grype
	./grype --fail-on=medium --only-fixed --config=.circleci/.anchore/grype.yaml -vv replicated/support-bundle:base