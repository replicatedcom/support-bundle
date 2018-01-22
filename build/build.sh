#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

if [ -z "${PKG}" ]; then
    echo "PKG must be set"
    exit 1
fi
if [ -z "${ARCH}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${VERSION}" ]; then
    echo "VERSION must be set"
    exit 1
fi
if [ -z "${SHA}" ]; then
    echo "SHA must be set"
    exit 1
fi
if [ -z "${BUILD_TIME}" ]; then
    echo "BUILD_TIME must be set"
    exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"

go install \
    -installsuffix "static" \
    -ldflags " \
    -X ${PKG}/pkg/version.version=${VERSION} \
    -X ${PKG}/pkg/version.gitSHA=${SHA} \
    -X ${PKG}/pkg/version.buildTime=${BUILD_TIME}" \
    ./cmd/...