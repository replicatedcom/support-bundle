FROM golang:1.9
WORKDIR /go/src/github.com/replicatedcom/support-bundle/
COPY . .
RUN make build

FROM debian:stretch-slim
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    sensible-utils=0.0.9+deb9u1 \
    ca-certificates \
    procps \
    systemd \
  && rm -rf /var/lib/apt/lists/*
COPY --from=0 /go/src/github.com/replicatedcom/support-bundle/bin/support-bundle /usr/local/bin/support-bundle
ENTRYPOINT [ "support-bundle" ]
CMD [ "generate" ]
