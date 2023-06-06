#!/bin/sh

set -o errexit
set -o nounset

go install github.com/onsi/ginkgo/ginkgo@latest

if [ -n "${DOCKER+x}" ]; then
    echo "Running e2e tests (docker enabled):"
    ginkgo -v -r -p --focus="docker" e2e/collect/core
    ginkgo -v -r -p --skip="swarm" e2e/collect/docker
    # ginkgo -v -r -p --focus="docker" e2e/collect/journald
elif [ -n "${SWARM+x}" ]; then
    echo "Running e2e tests (swarm enabled):"
    ginkgo -v -r -p --focus="swarm" e2e/collect/docker
elif [ -n "${RETRACED+x}" ]; then
    echo "Running e2e tests (retraced enabled):"
    ginkgo -v -r -p e2e/collect/retraced
else
    echo "Running e2e tests (core enabled):"
    ginkgo -v -r -p --skip="docker" e2e/collect/core
    # ginkgo -v -r -p --skip="docker" e2e/collect/journald
    ginkgo -v -r -p e2e/collect/supportbundle
    ginkgo -v -r -p e2e/collect/meta
fi
echo

echo "PASS"
