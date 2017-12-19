FROM alpine:3.2

#RUN apk update && \
  #apk add \
    #ca-certificates && \
  #rm -rf /var/cache/apk/*

COPY cron-server /bin

ENTRYPOINT ["/bin/cron-server"]
