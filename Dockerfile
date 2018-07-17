FROM golang:1.10

ENV PROJECTPATH=/go/src/github.com/replicatedcom/support-bundle

WORKDIR $PROJECTPATH

ADD Makefile .
RUN make build-deps dep-deps

CMD ["/bin/bash"]
