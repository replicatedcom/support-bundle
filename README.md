# Support bundle generator

This can be run in a Docker container or as a stand alone binary:

```bash
make support-bundle-generate
```

## Unit tests

```bash
make test
```

## Integration tests

```bash
make e2e-supportbundle-core e2e-supportbundle-docker
```

## Scanning image prior to release

```
make scan-base
```

## Releases

Releases are created on CircleCI when a tag is pushed.

```bash
git tag -a v0.1.0 -m "Initial release" && git push origin v0.1.0
```
