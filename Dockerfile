FROM golang:1.20

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    jq \
  && rm -rf /var/lib/apt/lists/*

ENV PROJECTPATH=/go/src/github.com/replicatedcom/support-bundle

WORKDIR $PROJECTPATH

ADD Makefile .
RUN make build-deps

CMD ["/bin/bash"]
