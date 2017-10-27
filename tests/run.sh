#!/bin/sh

docker run --rm \
    -v `pwd`:/go/src/github.com/replicatedcom/support-bundle \
    -w /go/src/github.com/replicatedcom/support-bundle/tests \
    golang:1.9-alpine \
    go test -v .