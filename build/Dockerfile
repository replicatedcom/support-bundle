FROM golang:1.9
WORKDIR /go/src/github.com/replicatedcom/support-bundle/
COPY . .
RUN make build

FROM debian:stretch-slim
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
  sensible-utils \
  ca-certificates \
  procps \
  systemd \
  iproute2 \
  && rm -rf /var/lib/apt/lists/*
COPY --from=0 /go/src/github.com/replicatedcom/support-bundle/bin/amd64/support-bundle /usr/local/bin/support-bundle
ENV IN_CONTAINER 1
LABEL "com.replicated.support-bundle"="true"
ENTRYPOINT [ "support-bundle" ]
CMD [ "generate" ]
