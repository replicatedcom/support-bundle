# Structure adapted from https://github.com/thockin/go-build-template

SHELL := /bin/bash
BIN := support-bundle
PKG := github.com/replicatedcom/support-bundle
VERSION := $(shell git describe --tags --always --dirty)
SHA := $(shell git log --pretty=format:'%H' -n 1)
ARCH ?= amd64

ifeq ($(shell uname), Darwin)
	BUILD_TIME := $(shell date -u +%FT%T)
else
	BUILD_TIME := $(shell date --rfc-3339=seconds | sed 's/ /T/')
endif

SRC_DIRS := cmd pkg
BUILD_IMAGE ?= golang:1.9-alpine
.PHONY: clean deps install run test build shell all

deps:
	go install

install:
	go install

generate:
	./bin/support-bundle generate
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

build-%:
	@$(MAKE) --no-print-directory ARCH=$* build

container-%:
	@$(MAKE) --no-print-directory ARCH=$* container

push-%:
	@$(MAKE) --no-print-directory ARCH=$* push

build: bin/$(ARCH)/$(BIN)

bin/$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
	    -u $$(id -u):$$(id -g)                                              \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
	    -v "$$(pwd)/bin/$(ARCH):/go/bin"                                    \
	    -v "$$(pwd)/bin/$(ARCH):/go/bin/$$(go env GOOS)_$(ARCH)"            \
	    -v "$$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static" \
	    -w /go/src/$(PKG)                                                   \
	    $(BUILD_IMAGE)                                                      \
	    /bin/sh -c "                                                        \
	        ARCH=$(ARCH)                                                    \
	        VERSION=$(VERSION)                                              \
	        PKG=$(PKG)                                                      \
			SHA=$(SHA)                                                      \
			BUILD_TIME=$(BUILD_TIME)                                        \
	        ./build/build.sh                                                \
	    "

# Example: make shell CMD="-c 'date > datefile'"
shell: build-dirs
	@echo "launching a shell in the containerized build environment"
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
	    -u $$(id -u):$$(id -g)                                              \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
	    -v "$$(pwd)/bin/$(ARCH):/go/bin"                                    \
	    -v "$$(pwd)/bin/$(ARCH):/go/bin/$$(go env GOOS)_$(ARCH)"            \
	    -v "$$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static" \
	    -w /go/src/$(PKG)                                                   \
	    $(BUILD_IMAGE)                                                      \
	    /bin/sh $(CMD)

test: build-dirs
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
	    -v "$$(pwd)/bin/$(ARCH):/go/bin"                                    \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
	    -v "$$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static" \
	    -w /go/src/$(PKG)                                                   \
	    $(BUILD_IMAGE)                                                      \
	    /bin/sh -c "                                                        \
	        ./build/test.sh $(SRC_DIRS)                                     \
	    "

e2e: build-dirs
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
	    -w /go/src/$(PKG)                                                   \
	    golang:1.9                                                          \
	    /bin/sh -c "                                                        \
	        ./build/e2e.sh $(SRC_DIRS)                                      \
	    "

e2e-docker: build-dirs
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
		--label com.replicated.support-bundle=true                          \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
	    -w /go/src/$(PKG)                                                   \
		golang:1.9                                                          \
	    /bin/sh -c "                                                        \
			DOCKER=1                                                        \
	        ./build/e2e.sh $(SRC_DIRS)                                      \
	    "

e2e-swarm: build-dirs
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
	    -w /go/src/$(PKG)                                                   \
		golang:1.9                                                          \
	    /bin/sh -c "                                                        \
			SWARM=1                                                         \
	        ./build/e2e.sh $(SRC_DIRS)                                      \
	    "

e2e-retraced: build-dirs
	@docker run                                                             \
	    -ti                                                                 \
	    --rm                                                                \
	    -v "$$(pwd)/.go:/go"                                                \
	    -v "$$(pwd):/go/src/$(PKG)"                                         \
		-v /var/run/docker.sock:/var/run/docker.sock                        \
	    -w /go/src/$(PKG)                                                   \
		golang:1.9                                                          \
	    /bin/sh -c "                                                        \
			RETRACED=1                                                      \
	        ./build/e2e.sh $(SRC_DIRS)                                      \
	    "

build-dirs:
	@mkdir -p bin/$(ARCH)
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/$(ARCH)

githooks:
	echo 'make integration-test' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	echo 'go fmt ./...' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

clean: container-clean bin-clean

container-clean:
	rm -rf .container-* .dockerfile-* .push-*

bin-clean:
	rm -rf .go bin

all: build test e2e e2e-docker
