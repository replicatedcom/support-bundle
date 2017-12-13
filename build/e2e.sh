#!/bin/sh

set -o errexit
set -o nounset

go get github.com/onsi/ginkgo/ginkgo

if [ -n "${DOCKER+x}" ]; then
    echo "Running e2e tests (docker enabled):"
    ginkgo -v -r -p --focus="docker container|journald.logs" e2e/ginkgo
else
    echo "Running e2e tests (docker disabled):"
    ginkgo -v -r -p --skip="docker container|retraced.events|journald.logs" e2e/ginkgo
fi
echo

echo "PASS"
