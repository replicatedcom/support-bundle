#!/bin/sh

docker build -t registry.replicated.com/library/support-bundle -f hack/Dockerfile .
