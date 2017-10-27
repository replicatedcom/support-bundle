#!/bin/sh

docker run --rm \
    -v `pwd`:/go/src/github.com/replicatedcom/support-bundle \
    -w /go/src/github.com/replicatedcom/support-bundle/tests \
    -v /var/run/docker.sock:/var/run/docker.sock \
    golang:1.9-alpine \
    go test -v .