FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY support-bundle /usr/local/bin/support-bundle
ENTRYPOINT [ "support-bundle" ]
CMD [ "generate" ]
