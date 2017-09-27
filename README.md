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

