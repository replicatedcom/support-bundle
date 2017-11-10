FROM golang:1.9
WORKDIR /go/src/github.com/replicatedcom/support-bundle/
COPY . .
RUN go build -o bin/support-bundle .

FROM debian:9
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    procps \
  && rm -rf /var/lib/apt/lists/*
COPY --from=0 /go/src/github.com/replicatedcom/support-bundle/bin/support-bundle /usr/local/bin/support-bundle
ENTRYPOINT [ "support-bundle" ]
CMD [ "generate" ]
