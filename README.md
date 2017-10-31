# support.io support bundle generator

This can be run in a Docker container or as a stand alone binary:

```bash
docker run \
  --pid=host \
  --net-host \
   -v /:/host:ro \
  -v`pwd`:/out \
  replicatedcom/support-bundle generate
```

### Unit tests
- run with `make test`

### Integration tests

- ginkgo suite at tests/ginkgo
- run with `make integration-test integration-test-docker`
