FROM golang:1.13

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    jq \
    \
    \
    libldap-2.4-2 \
  && apt-get clean \
  && apt-get autoremove -y \
  && rm -rf /var/lib/apt/lists/*

ENV PROJECTPATH=/go/src/github.com/replicatedcom/support-bundle

WORKDIR $PROJECTPATH

ADD Makefile .
RUN make build-deps

CMD ["/bin/bash"]
