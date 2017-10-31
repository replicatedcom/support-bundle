#!/bin/sh

docker run --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /dev:/host/dev \
    -v /proc:/host/proc:ro \
    -v /boot:/host/boot:ro \
    -v /lib/modules:/host/lib/modules:ro \
    -v /usr:/host/usr:ro \
    -v `pwd`:/out \
    -e IN_CONTAINER=true \
    replicated/support-bundle \
    generate --out /out/supportbundle.tar.gz
