#!/bin/sh

set -o errexit
set -o nounset

go get github.com/onsi/ginkgo/ginkgo

if [ -n "${DOCKER+x}" ]; then
    echo "Running e2e tests (docker enabled):"
    ginkgo -v -r -p --focus="docker" e2e/core
    ginkgo -v -r -p --skip="swarm" e2e/docker
    ginkgo -v -r -p --focus="docker" e2e/journald
elif [ -n "${SWARM+x}" ]; then
    echo "Running e2e tests (swarm enabled):"
    ginkgo -v -r -p --focus="swarm" e2e/docker
elif [ -n "${RETRACED+x}" ]; then
    echo "Running e2e tests (retraced enabled):"
    ginkgo -v -r -p e2e/retraced
else
    echo "Running e2e tests (core enabled):"
    ginkgo -v -r -p --skip="docker" e2e/core
    # ginkgo -v -r -p --skip="docker" e2e/journald
    ginkgo -v -r -p e2e/supportbundle
fi
echo

echo "PASS"
