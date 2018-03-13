FROM golang:1.9
WORKDIR /go/src/github.com/replicatedcom/support-bundle/
COPY . .
RUN make build

FROM debian:stretch-slim
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    procps \
    iproute2 \
    sensible-utils \
    systemd \
    libc6 \
    libc-bin \
    multiarch-support \
    libudev1 \
    libncurses5 \
    ncurses-base \
    ncurses-bin \
    libblkid1 \
    libfdisk1 \
    libmount1 \
    libsmartcols1 \
    libuuid1 \
    mount \
    util-linux \
  && rm -rf /var/lib/apt/lists/*
COPY --from=0 /go/src/github.com/replicatedcom/support-bundle/bin/amd64/support-bundle /usr/local/bin/support-bundle
ENV IN_CONTAINER 1
LABEL "com.replicated.support-bundle"="true"
ENTRYPOINT [ "support-bundle" ]
CMD [ "generate" ]
