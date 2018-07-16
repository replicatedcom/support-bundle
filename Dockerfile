FROM golang:1.10

ENV PROJECTPATH=/go/src/github.com/replicatedhq/troubleshoot

WORKDIR $PROJECTPATH

ADD Makefile .
RUN make build-deps dep-deps

CMD ["/bin/bash"]
