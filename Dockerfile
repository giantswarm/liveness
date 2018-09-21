FROM alpine:3.8

RUN apk add --update ca-certificates \
    && apk add --no-cache curl \
    && rm -rf /var/cache/apk/*

ADD ./liveness /liveness

ENTRYPOINT ["/liveness"]
