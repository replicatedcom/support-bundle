FROM golang:1.9-alpine
WORKDIR /go/src/github.com/replicatedcom/support-bundle/
COPY . .
RUN go build -o bin/support-bundle .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/replicatedcom/support-bundle/bin/support-bundle /usr/local/bin/support-bundle
CMD ["support-bundle"]  
