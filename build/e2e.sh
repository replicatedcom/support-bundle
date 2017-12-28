#!/bin/sh

set -o errexit
set -o nounset

go get github.com/onsi/ginkgo/ginkgo

if [ -n "${DOCKER+x}" ]; then
    echo "Running e2e tests (docker enabled):"
    ginkgo -v -r -p --focus="docker|journald" --skip="swarm" e2e/
elif [ -n "${SWARM+x}" ]; then
    echo "Running e2e tests (swarm enabled):"
    ginkgo -v -r -p --focus="docker swarm" e2e/
elif [ -n "${RETRACED+x}" ]; then
    echo "Running e2e tests (retraced enabled):"
    ginkgo -v -r -p --focus="retraced" e2e/
else
    echo "Running e2e tests (core):"
    ginkgo -v -r -p --skip="docker|journald|retraced" e2e/
fi
echo

echo "PASS"
