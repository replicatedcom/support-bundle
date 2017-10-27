#!/bin/sh

docker run --rm \
    --pid=host \
    --net=host \
    -v `pwd`:/out \
    replicated/support-bundle \
    support-bundle generate --out /out/supportbundle.tar.gz
