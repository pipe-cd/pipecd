FROM alpine:3.10.1

RUN \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    rm -f /var/cache/apk/*
