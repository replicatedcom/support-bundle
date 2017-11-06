# Support bundle generator

This can be run in a Docker container or as a stand alone binary:

```bash
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
```

## Unit tests

```bash
make test
```

## Integration tests

```bash
make integration-test integration-test-docker
```

## Releases

Releases are created on CircleCI when a tag is pushed.

```bash
git tag -a v0.1.0 -m "First release" && git push origin v0.1.0
```
