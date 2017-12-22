#!/bin/sh

set -o errexit
set -o nounset

go get github.com/onsi/ginkgo/ginkgo

if [ -n "${DOCKER+x}" ]; then
    echo "Running e2e tests (docker enabled):"
    ginkgo -v -r -p --focus="docker|journald" --skip="swarm" e2e/ginkgo
elif [ -n "${SWARM+x}" ]; then
    echo "Running e2e tests (swarm enabled):"
    ginkgo -v -r -p --focus="swarm" e2e/ginkgo
else
    echo "Running e2e tests (docker disabled):"
    ginkgo -v -r -p --skip="docker|swarm|retraced|journald" e2e/ginkgo
fi
echo

echo "PASS"
