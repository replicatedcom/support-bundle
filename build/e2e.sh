#!/bin/sh

set -o errexit
set -o nounset

go get github.com/onsi/ginkgo/ginkgo

if [ -n "${DOCKER+x}" ]; then
    echo "Running e2e tests (docker enabled):"
    docker pull ubuntu:latest
    ginkgo -v -r -p --focus="docker container|journald.logs" tests/ginkgo
else
    echo "Running e2e tests (docker disabled):"
    ginkgo -v -r -p --skip="docker container|retraced.events|journald.logs" tests/ginkgo
fi
echo

echo "PASS"
